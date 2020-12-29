package vault

import (
	"github.com/polarctos/situ-vault/pkg/testdata"
	"github.com/polarctos/situ-vault/pkg/vault/mode"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_roundTrip(t *testing.T) {
	password := string(testdata.RandomPassword(16))
	cleartext := string(testdata.RandomDataBase64(500))
	modeText := mode.Defaults().Conservative.Text()

	resultEncrypted, err := Encrypt(cleartext, password, modeText)
	assert.Nil(t, err)
	assert.Contains(t, resultEncrypted, "SITU_VAULT")
	assert.NotContains(t, resultEncrypted, cleartext)

	resultDecrypted, modeTextResult, err := Decrypt(resultEncrypted, password)
	assert.Nil(t, err)
	assert.Contains(t, resultDecrypted, cleartext)
	assert.Contains(t, modeTextResult, modeText)
}
