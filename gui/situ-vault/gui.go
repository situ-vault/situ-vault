package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"

	"github.com/situ-vault/situ-vault/pkg/vault"
	"github.com/situ-vault/situ-vault/pkg/vault/vaultmode"
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
	passwordFile   func()
	passwordNotify *ToolbarNotification
	input          *widget.Entry
	inputPaste     func()
	inputFile      func()
	modes          *widget.RadioGroup
	modesAddCustom *ToolbarLabeledAction
	modesDialog    *customModesDialog
	action         func()
	clear          func()
	output         *widget.Entry
	outputCut      func()
	outputCopy     func()
	outputFile     func()
	outputNotify   *ToolbarNotification
	clearClipboard func()
	clearNotify    *ToolbarNotification
	refreshObjects []interface{ fyne.CanvasObject }
}

func newUi(w fyne.Window, model *model, action func(), refresh func(), showError func(error), getClipboard func() fyne.Clipboard) *ui {
	u := &ui{}
	u.refreshObjects = make([]interface{ fyne.CanvasObject }, 0, 2)
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
	refreshAll := func() {
		refresh()
		for _, refreshable := range u.refreshObjects {
			refreshable.Refresh()
		}
	}
	updateUiFromModel := func() {
		u.password.Text = model.password
		u.input.Text = model.input
		syncModes()
		u.modes.Selected = model.mode
		u.output.Text = model.output
		refreshAll()
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
		u.password.Password = true
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
	u.outputNotify = NewToolbarNotification()
	u.outputCopy = func() {
		getClipboard().SetContent(model.output)
		u.outputNotify.ShowNotification("Copied to clipboard.", time.Second, refreshAll)
	}
	u.outputFile = func() {
		if model.output == "" {
			showError(errors.New("no value to save yet"))
			return
		}
		callback := func(file fyne.URIWriteCloser, err error) {
			if err == nil && file != nil {
				_, err := file.Write([]byte(model.output))
				if err != nil {
					showError(err)
				}
			}
		}
		fileSaveDialog := dialog.NewFileSave(callback, w)
		listUri, err := wdListUri()
		if err != nil {
			showError(err)
		}
		fileSaveDialog.SetLocation(listUri)
		fileSaveDialog.Resize(fyne.NewSize(700, 600))
		fileSaveDialog.Show()
	}

	u.passwordNotify = NewToolbarNotification()
	u.passwordCopy = func() {
		updateModelFromUi()
		getClipboard().SetContent(model.password)
		u.passwordNotify.ShowNotification("Password copied.", time.Second, refreshAll)
	}
	u.passwordPaste = func() {
		updateModelFromUi()
		model.password = getClipboard().Content()
		updateUiFromModel()
	}
	u.passwordFile = func() {
		filePathViaDialog(w, showError, func(filePath string) {
			updateModelFromUi()
			model.password = filePath
			u.password.Password = false
			updateUiFromModel()
		})
	}

	u.inputPaste = func() {
		updateModelFromUi()
		model.input = getClipboard().Content()
		updateUiFromModel()
	}
	u.inputFile = func() {
		filePathViaDialog(w, showError, func(filePath string) {
			updateModelFromUi()
			model.input = filePath
			updateUiFromModel()
		})
	}

	u.clearNotify = NewToolbarNotification()
	u.clearClipboard = func() {
		getClipboard().SetContent("")
		u.clearNotify.ShowNotification("Clipboard cleared.", time.Second, refreshAll)
	}

	u.modesAddCustom = NewToolbarLabeledAction(theme.ColorPaletteIcon(), "Add Custom Mode", func() {
		showCustomModeDialog(w, u, refresh)
	})

	return u
}

func filePathViaDialog(w fyne.Window, showError func(error), filePathCallback func(string)) {
	callback := func(file fyne.URIReadCloser, err error) {
		if err == nil && file != nil {
			uri := file.URI()
			filePathCallback(uri.String())
		}
	}
	fileOpenDialog := dialog.NewFileOpen(callback, w)
	listUri, err := wdListUri()
	if err != nil {
		showError(err)
	}
	fileOpenDialog.SetLocation(listUri)
	fileOpenDialog.Resize(fyne.NewSize(700, 600))
	fileOpenDialog.Show()
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
	return newUi(w, model, action, refresh, showError, getClipboard)
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
	return newUi(w, model, action, refresh, showError, getClipboard)
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
	w.Resize(fyne.NewSize(1000, 900))

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
	content := container.New(
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
					NewToolbarLabeledAction(theme.FolderOpenIcon(), "File", ui.passwordFile),
					ui.passwordNotify,
				),
			)},
			{Text: inputName + ":", Widget: container.NewVBox(
				scrollableMultilineEntry(ui.input, ui),
				widget.NewToolbar(
					NewToolbarLabeledAction(theme.ContentPasteIcon(), "Paste", ui.inputPaste),
					NewToolbarLabeledAction(theme.FolderOpenIcon(), "File", ui.inputFile),
				),
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
		NewToolbarLabeledAction(theme.FolderIcon(), "Save", ui.outputFile),
		ui.outputNotify,
	)
	result := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text: string(op) + "ed:",
				Widget: container.NewVBox(
					scrollableMultilineEntry(ui.output, ui),
					toolbar,
				)},
		},
		OnCancel:   ui.clearClipboard,
		CancelText: "Clear Clipboard",
	}

	// vaultmode is either before or after the separator:
	modeWidget := container.NewVBox(ui.modes)
	modeFormItem := widget.NewFormItem("Mode:", modeWidget)
	switch op {
	case Encrypt:
		modeWidget.Add(widget.NewToolbar(ui.modesAddCustom))
		form.Items = append(form.Items, modeFormItem)
	case Decrypt:
		result.Items = append([]*widget.FormItem{modeFormItem}, result.Items...)
	}

	space := widget.NewLabel("")
	separator := widget.NewSeparator()

	clipboardNotifyBox := container.NewHBox(layout.NewSpacer(), ui.clearNotify.ToolbarObject())

	return container.NewVBox(
		space,
		infoText,
		form,
		space,
		separator,
		space,
		result,
		clipboardNotifyBox,
	)
}

func scrollableMultilineEntry(e *widget.Entry, u *ui) fyne.CanvasObject {
	scrollable := container.NewVScroll(e)
	scrollable.SetMinSize(e.MinSize().Add(fyne.NewSize(0, 50)))
	// workaround as fyne does not recursively refresh the contents of Scrolls:
	u.refreshObjects = append(u.refreshObjects, e)
	return scrollable
}
