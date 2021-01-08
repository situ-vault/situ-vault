package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/situ-vault/situ-vault/pkg/testdata"
	"github.com/situ-vault/situ-vault/pkg/vault/vaultmode"
)

func Test_roundTrip_Conservative(t *testing.T) {
	testRound(t, vaultmode.Defaults().Conservative.Text())
}
func Test_roundTrip_Modern(t *testing.T) {
	testRound(t, vaultmode.Defaults().Modern.Text())
}
func Test_roundTrip_Secretbox(t *testing.T) {
	testRound(t, vaultmode.Defaults().Secretbox.Text())
}
func Test_roundTrip_XChaCha(t *testing.T) {
	testRound(t, vaultmode.Defaults().XChaCha.Text())
}

func testRound(t *testing.T, modeText string) {
	password := string(testdata.RandomPassword(16))
	cleartext := string(testdata.RandomDataBase64(500))

	resultEncrypted, err := Encrypt(cleartext, password, modeText)
	assert.Nil(t, err)
	assert.Contains(t, resultEncrypted, "SITU_VAULT")
	assert.NotContains(t, resultEncrypted, cleartext)

	resultDecrypted, modeTextResult, err := Decrypt(resultEncrypted, password)
	assert.Nil(t, err)
	assert.Contains(t, resultDecrypted, cleartext)
	assert.Contains(t, modeTextResult, modeText)
}
