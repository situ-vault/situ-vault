package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_envelope(t *testing.T) {
	salt := "salt"
	ciphertext := "ciphertext"

	envelope := buildEnvelope(salt, ciphertext)
	assert.EqualValues(t, "SITU_VAULT_V1##AES256_GCM_PBKDF2_SHA256_ITER10K_SALT8_BASE32##salt##ciphertext", envelope)

	openedSalt, openedCipherText, _ := openEnvelope(envelope)
	assert.EqualValues(t, salt, openedSalt)
	assert.EqualValues(t, ciphertext, openedCipherText)
}
