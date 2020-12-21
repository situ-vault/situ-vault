package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/polarctos/situ-vault/pkg/vault"
	"image/color"
)

type operation string

const (
	Encrypt operation = "Encrypt"
	Decrypt operation = "Decrypt"
)

var modes = []string{"AES256_GCM_PBKDF2_SHA256_ITER10K_SALT8_BASE32"}

type model struct {
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
	action         func()
	clear          func()
	output         *widget.Entry
	outputCut      func()
	outputCopy     func()
	clearClipboard func()
}

func newUi(model *model, action func(), refresh func(), getClipboard func() fyne.Clipboard) *ui {
	u := &ui{}
	u.password = widget.NewPasswordEntry()
	u.input = widget.NewMultiLineEntry()
	u.input.Wrapping = fyne.TextWrapBreak
	u.modes = widget.NewRadioGroup(modes, func(string) {})
	u.modes.SetSelected(modes[0])

	// there is no data binding in fyne yet, thus manually:
	updateModelFromUi := func() {
		model.password = u.password.Text
		model.input = u.input.Text
		model.mode = u.modes.Selected
	}
	updateUiFromModel := func() {
		u.password.Text = model.password
		u.input.Text = model.input
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

	return u
}

func newEncryptUi(refresh func(), showError func(error), getClipboard func() fyne.Clipboard) *ui {
	model := newModel()
	action := func() {
		result, err := vault.Encrypt(model.input, model.password)
		if err != nil {
			showError(err)
		} else {
			model.output = result
		}
	}
	return newUi(model, action, refresh, getClipboard)
}

func newDecryptUi(refresh func(), showError func(error), getClipboard func() fyne.Clipboard) *ui {
	model := newModel()
	action := func() {
		result, err := vault.Decrypt(model.input, model.password)
		if err != nil {
			showError(err)
		} else {
			model.mode = modes[0] // FIXME needed next to decryption result if multiple modes exist
			model.output = result
		}
	}
	return newUi(model, action, refresh, getClipboard)
}

func newModel() *model {
	return &model{
		password: "",
		input:    "",
		mode:     modes[0],
	}
}

func main() {
	a := app.New()

	w := a.NewWindow("situ-vault")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(800, 500))

	showError := func(err error) {
		dialog.ShowError(err, w)
	}

	getClipboard := func() fyne.Clipboard {
		return w.Clipboard() // clipboard only available when window was shown
	}

	var encryptTab *fyne.Container
	encryptUi := newEncryptUi(func() { encryptTab.Refresh() }, showError, getClipboard)
	encryptTab = uiTabDesign(encryptUi, Encrypt)

	var decryptTab *fyne.Container
	decryptUi := newDecryptUi(func() { decryptTab.Refresh() }, showError, getClipboard)
	decryptTab = uiTabDesign(decryptUi, Decrypt)

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
	w.ShowAndRun()
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

	// mode is either before or after the separator:
	mode := widget.NewFormItem("Mode:", ui.modes)
	switch op {
	case Encrypt:
		form.Items = append(form.Items, mode)
	case Decrypt:
		result.Items = append([]*widget.FormItem{mode}, result.Items...)
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
