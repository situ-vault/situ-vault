package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var salt = newSalt(SaltLength8)

func Test_aesgcm(t *testing.T) {
	key := deriveKey([]byte("test-pw"), salt)
	data := []byte("test-data")

	ciphertext, err := encrypt(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)
	cleartext, err := decrypt(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)

	wrongKey := deriveKey([]byte("wrong-pw"), salt) // wrong pw, same salt
	cleartext, err = decrypt(ciphertext, wrongKey)
	assert.NotNil(t, err, "wrong password should not decrypt")
}

func Test_secretbox(t *testing.T) {
	key := deriveKey([]byte("test-pw"), salt)
	data := []byte("test-data")

	ciphertext, err := encryptSecretbox(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)
	cleartext, err := decryptSecretbox(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)

	wrongKey := deriveKey([]byte("wrong-pw"), salt) // wrong pw, same salt
	cleartext, err = decryptSecretbox(ciphertext, wrongKey)
	assert.NotNil(t, err, "wrong password should not decrypt")
}

func Test_XChaCha20Poly1305(t *testing.T) {
	key := deriveKey([]byte("test-pw"), salt)
	data := []byte("test-data")

	ciphertext, err := encryptXChaPo(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)
	cleartext, err := decryptXChaPo(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)

	wrongKey := deriveKey([]byte("wrong-pw"), salt) // wrong pw, same salt
	cleartext, err = decryptXChaPo(ciphertext, wrongKey)
	assert.NotNil(t, err, "wrong password should not decrypt")
}
