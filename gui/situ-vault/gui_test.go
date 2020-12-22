package main

import (
	"fyne.io/fyne/test"
	"github.com/polarctos/situ-vault/pkg/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

var predefined = testdata.PredefinedDecrypt()

func newTestApp() *experience {
	exp := newExperience()
	exp.loadUi(test.NewApp())
	return exp
}

func Test_encrypt(t *testing.T) {
	ui := newTestApp().ui[Encrypt]

	test.Type(ui.password, testdata.RandomPassword(5))
	test.Type(ui.input, testdata.RandomDataBase64(200))
	assert.Equal(t, "", ui.output.Text, "No output before encryption")

	// unfortunately buttons inside a form are private, thus not directly tappable in the test
	ui.action()
	assert.Contains(t, ui.output.Text, "SITU_VAULT", "Encrypted output after interaction")
}

func Test_decrypt(t *testing.T) {
	ui := newTestApp().ui[Decrypt]

	test.Type(ui.password, predefined.Password)
	test.Type(ui.input, predefined.Ciphertext)
	assert.Equal(t, "", ui.output.Text, "No output before decryption")

	// unfortunately buttons inside a form are private, thus not directly tappable in the test
	ui.action()
	assert.Contains(t, ui.output.Text, predefined.Cleartext, "Decrypted output after interaction")
}
