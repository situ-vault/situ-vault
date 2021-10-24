package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/situ-vault/situ-vault/pkg/vault/vaultmode"
)

type customModesDialog struct {
	dialog      dialog.Dialog
	callback    func(ok bool)
	modeBuilder *modeBuilder
}

type modeBuilder struct {
	construct *widget.RadioGroup
	kdf       *widget.RadioGroup
	salt      *widget.RadioGroup
	encoding  *widget.RadioGroup
	linebreak *widget.RadioGroup
}

func showCustomModeDialog(w fyne.Window, u *ui, refresh func()) {
	cmd := customModesDialog{}
	u.modesDialog = &cmd
	modeBuilder := &modeBuilder{
		construct: widget.NewRadioGroup(vaultmode.Constructs.AllValues(), func(string) {}),
		kdf:       widget.NewRadioGroup(vaultmode.KeyDerivationFunctions.AllValues(), func(string) {}),
		salt:      widget.NewRadioGroup(vaultmode.Salts.AllValues(), func(string) {}),
		encoding:  widget.NewRadioGroup(vaultmode.Encodings.AllValues(), func(string) {}),
		linebreak: widget.NewRadioGroup(vaultmode.Linebreaks.AllValues(), func(string) {}),
	}
	// pre-select the first element, thus it is always a valid setup:
	modeBuilder.construct.SetSelected(vaultmode.Constructs.AllValues()[0])
	modeBuilder.kdf.SetSelected(vaultmode.KeyDerivationFunctions.AllValues()[0])
	modeBuilder.salt.SetSelected(vaultmode.Salts.AllValues()[0])
	modeBuilder.encoding.SetSelected(vaultmode.Encodings.AllValues()[0])
	modeBuilder.linebreak.SetSelected(vaultmode.Linebreaks.AllValues()[0])
	cmd.modeBuilder = modeBuilder
	cmd.callback = func(ok bool) {
		if ok {
			customMode := vaultmode.Mode{
				Construct: vaultmode.Construct(modeBuilder.construct.Selected),
				Kdf:       vaultmode.KeyDerivationFunction(modeBuilder.kdf.Selected),
				Salt:      vaultmode.Salt(modeBuilder.salt.Selected),
				Encoding:  vaultmode.Encoding(modeBuilder.encoding.Selected),
				Linebreak: vaultmode.Linebreak(modeBuilder.linebreak.Selected),
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
			{Text: "Linebreak:", Widget: modeBuilder.linebreak},
		},
	}
	cmd.dialog = dialog.NewCustomConfirm("Build Custom Mode", "Add", "Cancel", content, cmd.callback, w)
	cmd.dialog.Show()
}
