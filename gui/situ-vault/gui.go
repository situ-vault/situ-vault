package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/polarctos/situ-vault/pkg/vault"
)

type ui struct {
	password   *widget.Entry
	input      *widget.Entry
	modes      *widget.RadioGroup
	action     func()
	clear      func()
	output     *widget.Entry
	outputCut  func()
	outputCopy func()
}

func newEncryptUi(showError func(error), getClipboard func() fyne.Clipboard) *ui {
	u := &ui{}
	u.password = widget.NewPasswordEntry()
	u.input = widget.NewMultiLineEntry()

	modes := []string{"AES256_GCM_PBKDF2_SHA256_ITER10K_SALT8_BASE32"}
	u.modes = widget.NewRadioGroup(modes, func(string) {})
	u.modes.SetSelected(modes[0])

	u.action = func() {
		result, err := vault.Encrypt(u.input.Text, u.password.Text)
		if err != nil {
			showError(err)
		} else {
			u.output.SetText(result)
		}
	}
	u.clear = func() {
		u.password.SetText("")
		u.input.SetText("")
		u.output.SetText("")
	}

	u.output = widget.NewMultiLineEntry()
	u.output.Wrapping = fyne.TextWrapBreak
	u.output.Disable()
	u.outputCut = func() {
		getClipboard().SetContent(u.output.Text)
		u.output.SetText("")
	}
	u.outputCopy = func() {
		getClipboard().SetContent(u.output.Text)
	}

	return u
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

	u := newEncryptUi(showError, getClipboard)

	infoText := widget.NewLabel("Encrypts text data using the selected algorithm.")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Password:", Widget: u.password},
			{Text: "Cleartext:", Widget: u.input},
			{Text: "Mode:", Widget: u.modes},
		},
		OnSubmit:   u.action,
		OnCancel:   u.clear,
		SubmitText: "Encrypt",
		CancelText: "Clear",
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentCutIcon(), u.outputCut),
		widget.NewToolbarAction(theme.ContentCopyIcon(), u.outputCopy),
	)

	result := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Encrypted:", Widget: container.NewVBox(u.output, toolbar)},
		},
	}

	separator := widget.NewSeparator()

	boxEncrypt := container.NewVBox(
		infoText,
		form,
		separator,
		result,
	)
	boxDecrypt := container.NewCenter(widget.NewLabel("TODO!"))
	appTabs := container.NewAppTabs(
		container.NewTabItem("Encrypt", boxEncrypt),
		container.NewTabItem("Decrypt", boxDecrypt),
	)

	w.SetContent(appTabs)
	w.ShowAndRun()
}
