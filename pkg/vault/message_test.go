package vault

import (
	"github.com/polarctos/situ-vault/pkg/vault/mode"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_message(t *testing.T) {
	salt := "salt"
	ciphertext := "ciphertext"

	m := Message{
		Prefix:     prefix,
		Mode:       mode.Defaults().Conservative,
		Salt:       salt,
		Ciphertext: ciphertext,
	}

	messageText := m.Text()
	assert.EqualValues(t, "SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##salt##ciphertext", messageText)

	openedMessage, err := NewMessage(messageText)
	assert.Nil(t, err)
	assert.EqualValues(t, prefix, openedMessage.Prefix)
	assert.EqualValues(t, mode.Defaults().Conservative, openedMessage.Mode)
	assert.EqualValues(t, salt, openedMessage.Salt)
	assert.EqualValues(t, ciphertext, openedMessage.Ciphertext)
}
