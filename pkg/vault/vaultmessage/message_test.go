package vaultmessage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/polarctos/situ-vault/pkg/vault/vaultmode"
)

func Test_message(t *testing.T) {
	salt := "salt"
	ciphertext := "ciphertext"

	m := Message{
		Prefix:     VaultPrefix,
		Mode:       vaultmode.Defaults().Conservative,
		Salt:       salt,
		Ciphertext: ciphertext,
	}

	messageText := m.Text()
	assert.EqualValues(t, "SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##salt##ciphertext", messageText)

	openedMessage, err := NewMessage(messageText)
	assert.Nil(t, err)
	assert.EqualValues(t, VaultPrefix, openedMessage.Prefix)
	assert.EqualValues(t, vaultmode.Defaults().Conservative, openedMessage.Mode)
	assert.EqualValues(t, salt, openedMessage.Salt)
	assert.EqualValues(t, ciphertext, openedMessage.Ciphertext)
}
