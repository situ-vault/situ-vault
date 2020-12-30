package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/polarctos/situ-vault/pkg/vault"
	"github.com/polarctos/situ-vault/pkg/vault/vaultmode"
)

type operation string

const (
	Encrypt operation = "Encrypt"
	Decrypt operation = "Decrypt"
)

var modes = []string{
	vaultmode.Defaults().Conservative.Text(),
	vaultmode.Defaults().Modern.Text(),
	vaultmode.Defaults().Secretbox.Text(),
	vaultmode.Defaults().XChaCha.Text(),
}

type experience struct {
	model map[operation]*model
	ui    map[operation]*ui
}

func newExperience() *experience {
	exp := &experience{
		model: make(map[operation]*model),
		ui:    make(map[operation]*ui),
	}
	return exp
}

type model struct {
	op       operation
	password string
	input    string
	mode     string
	output   string
}

type ui struct {
	password       *widget.Entry
	passwordPaste  func()
	passwordCopy   func()
	input          *widget.Entry
	inputPaste     func()
	modes          *widget.RadioGroup
	modesAddCustom func()
	action         func()
	clear          func()
	output         *widget.Entry
	outputCut      func()
	outputCopy     func()
	clearClipboard func()
}

type customModeBuilder struct {
	construct *widget.RadioGroup
	kdf       *widget.RadioGroup
	salt      *widget.RadioGroup
	encoding  *widget.RadioGroup
}

func newUi(w fyne.Window, model *model, action func(), refresh func(), getClipboard func() fyne.Clipboard) *ui {
	u := &ui{}
	u.password = widget.NewPasswordEntry()
	u.input = widget.NewMultiLineEntry()
	u.input.Wrapping = fyne.TextWrapBreak
	u.modes = widget.NewRadioGroup(modes, func(string) {})
	u.modes.SetSelected(modes[0])
	if model.op == Decrypt {
		u.modes.Disable()
	}

	// there is no data binding in fyne yet, thus manually:
	syncModes := func() {
		for _, mode := range modes {
			u.modes.Append(mode)
		}
	}
	updateModelFromUi := func() {
		model.password = u.password.Text
		model.input = u.input.Text
		model.mode = u.modes.Selected
	}
	updateUiFromModel := func() {
		u.password.Text = model.password
		u.input.Text = model.input
		syncModes()
		u.modes.Selected = model.mode
		u.output.Text = model.output
		refresh()
	}

	u.action = func() {
		updateModelFromUi()
		action()
		updateUiFromModel()
	}
	u.clear = func() {
		model.password = ""
		model.input = ""
		model.mode = modes[0]
		model.output = ""
		updateUiFromModel()
	}

	u.output = widget.NewMultiLineEntry()
	u.output.Wrapping = fyne.TextWrapBreak
	u.output.Disable()
	u.outputCut = func() {
		getClipboard().SetContent(model.output)
		model.output = ""
		updateUiFromModel()

	}
	u.outputCopy = func() {
		getClipboard().SetContent(model.output)
	}

	u.passwordCopy = func() {
		updateModelFromUi()
		getClipboard().SetContent(model.password)
	}
	u.passwordPaste = func() {
		updateModelFromUi()
		model.password = getClipboard().Content()
		updateUiFromModel()
	}
	u.inputPaste = func() {
		updateModelFromUi()
		model.input = getClipboard().Content()
		updateUiFromModel()
	}

	u.clearClipboard = func() {
		getClipboard().SetContent("")
	}

	modeBuilder := customModeBuilder{
		construct: widget.NewRadioGroup(vaultmode.Constructs.AllValues(), func(string) {}),
		kdf:       widget.NewRadioGroup(vaultmode.KeyDerivationFunctions.AllValues(), func(string) {}),
		salt:      widget.NewRadioGroup(vaultmode.Salts.AllValues(), func(string) {}),
		encoding:  widget.NewRadioGroup(vaultmode.Encodings.AllValues(), func(string) {}),
	}
	// pre-select the first element, thus it is always a valid setup:
	modeBuilder.construct.SetSelected(vaultmode.Constructs.AllValues()[0])
	modeBuilder.kdf.SetSelected(vaultmode.KeyDerivationFunctions.AllValues()[0])
	modeBuilder.salt.SetSelected(vaultmode.Salts.AllValues()[0])
	modeBuilder.encoding.SetSelected(vaultmode.Encodings.AllValues()[0])
	callback := func(ok bool) {
		if ok {
			customMode := vaultmode.Mode{
				Construct: vaultmode.Construct(modeBuilder.construct.Selected),
				Kdf:       vaultmode.KeyDerivationFunction(modeBuilder.kdf.Selected),
				Salt:      vaultmode.Salt(modeBuilder.salt.Selected),
				Encoding:  vaultmode.Encoding(modeBuilder.encoding.Selected),
			}
			customModeText := customMode.Text()
			u.modes.Append(customModeText)
			u.modes.SetSelected(customModeText)
			refresh()
		}
	}
	content := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Construct:", Widget: modeBuilder.construct},
			{Text: "Key Derivation Function:", Widget: modeBuilder.kdf},
			{Text: "Salt:", Widget: modeBuilder.salt},
			{Text: "Encoding:", Widget: modeBuilder.encoding},
		},
	}
	u.modesAddCustom = func() {
		dialog.ShowCustomConfirm("Build Custom Mode", "Add", "Cancel", content, callback, w)
	}

	return u
}

func newEncryptUi(w fyne.Window, model *model, refresh func(), showError func(error), getClipboard func() fyne.Clipboard) *ui {
	action := func() {
		result, err := vault.Encrypt(model.input, model.password, model.mode)
		if err != nil {
			showError(err)
		} else {
			model.output = result
		}
	}
	return newUi(w, model, action, refresh, getClipboard)
}

func newDecryptUi(w fyne.Window, model *model, refresh func(), showError func(error), getClipboard func() fyne.Clipboard) *ui {
	action := func() {
		result, modeText, err := vault.Decrypt(model.input, model.password)
		if err != nil {
			showError(err)
		} else {
			modes = append(modes, modeText)
			model.mode = modeText
			model.output = result
		}
	}
	return newUi(w, model, action, refresh, getClipboard)
}

func newModel(op operation) *model {
	return &model{
		op:       op,
		password: "",
		input:    "",
		mode:     modes[0],
	}
}

func (exp *experience) loadUi(application fyne.App) {

	w := application.NewWindow("situ-vault")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(800, 500))

	showError := func(err error) {
		dialog.ShowError(err, w)
	}

	getClipboard := func() fyne.Clipboard {
		return w.Clipboard() // clipboard only available when window was shown
	}

	var encryptTab *fyne.Container
	encryptModel := newModel(Encrypt)
	encryptUi := newEncryptUi(w, encryptModel, func() { encryptTab.Refresh() }, showError, getClipboard)
	encryptTab = uiTabDesign(encryptUi, Encrypt)
	exp.model[Encrypt] = encryptModel
	exp.ui[Encrypt] = encryptUi

	var decryptTab *fyne.Container
	decryptModel := newModel(Decrypt)
	decryptUi := newDecryptUi(w, decryptModel, func() { decryptTab.Refresh() }, showError, getClipboard)
	decryptTab = uiTabDesign(decryptUi, Decrypt)
	exp.model[Decrypt] = decryptModel
	exp.ui[Decrypt] = decryptUi

	appTabs := container.NewAppTabs(
		container.NewTabItem("Encrypt", encryptTab),
		container.NewTabItem("Decrypt", decryptTab),
	)

	// workaround to give it a bit more space to breathe
	space := canvas.NewRectangle(color.White)
	space.SetMinSize(fyne.NewSize(20, 20))
	content := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(space, space, space, space),
		appTabs,
	)

	w.SetContent(content)
	w.Show()
}

func uiTabDesign(ui *ui, op operation) *fyne.Container {
	infoText := widget.NewLabel(string(op) + "s text data using the selected algorithm.\n")

	var inputName string
	switch op {
	case Encrypt:
		inputName = "Cleartext"
	case Decrypt:
		inputName = "Ciphertext"
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Password:", Widget: container.NewVBox(
				ui.password,
				widget.NewToolbar(
					NewToolbarLabeledAction(theme.ContentPasteIcon(), "Paste", ui.passwordPaste),
					NewToolbarLabeledAction(theme.ContentCopyIcon(), "Copy", ui.passwordCopy),
				),
			)},
			{Text: inputName + ":", Widget: container.NewVBox(
				ui.input,
				widget.NewToolbar(
					NewToolbarLabeledAction(theme.ContentPasteIcon(), "Paste", ui.inputPaste)),
			)},
		},
		OnSubmit:   ui.action,
		OnCancel:   ui.clear,
		SubmitText: string(op),
		CancelText: "Clear",
	}
	toolbar := widget.NewToolbar(
		NewToolbarLabeledAction(theme.ContentCutIcon(), "Cut", ui.outputCut),
		NewToolbarLabeledAction(theme.ContentCopyIcon(), "Copy", ui.outputCopy),
	)
	result := &widget.Form{
		Items: []*widget.FormItem{
			{Text: string(op) + "ed:", Widget: container.NewVBox(ui.output, toolbar)},
		},
		OnCancel:   ui.clearClipboard,
		CancelText: "Clear Clipboard",
	}

	// vaultmode is either before or after the separator:
	modeWidget := container.NewVBox(ui.modes)
	modeFormItem := widget.NewFormItem("Mode:", modeWidget)
	switch op {
	case Encrypt:
		modeWidget.Add(widget.NewToolbar(
			NewToolbarLabeledAction(theme.ColorPaletteIcon(), "Add Custom Mode", ui.modesAddCustom)))
		form.Items = append(form.Items, modeFormItem)
	case Decrypt:
		result.Items = append([]*widget.FormItem{modeFormItem}, result.Items...)
	}

	space := widget.NewLabel("")
	separator := widget.NewSeparator()

	return container.NewVBox(
		space,
		infoText,
		form,
		space,
		separator,
		space,
		result,
	)
}
