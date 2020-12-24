package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_envelope(t *testing.T) {
	salt := "salt"
	ciphertext := "ciphertext"

	envelope := buildEnvelope(salt, ciphertext)
	assert.EqualValues(t, "SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##salt##ciphertext", envelope)

	openedSalt, openedCipherText, err := openEnvelope(envelope)
	assert.Nil(t, err)
	assert.EqualValues(t, salt, openedSalt)
	assert.EqualValues(t, ciphertext, openedCipherText)
}
