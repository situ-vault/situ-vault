package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_aesgcm(t *testing.T) {
	key := deriveKey([]byte("test-pw"), newSalt())
	data := []byte("test-data")

	ciphertext, err := encrypt(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)
	cleartext, err := decrypt(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)
}

func Test_secretbox(t *testing.T) {
	key := deriveKey([]byte("test-pw"), newSalt())
	data := []byte("test-data")

	secretbox, err := encryptSecretbox(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, secretbox, data)
	cleartext, err := decryptSecretbox(secretbox, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)
}

func Test_XChaCha20Poly1305(t *testing.T) {
	key := deriveKey([]byte("test-pw"), newSalt())
	data := []byte("test-data")

	ciphertext, err := encryptXChaPo(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)
	cleartext, err := decryptXChaPo(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)
}
