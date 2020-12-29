package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var salt = NewSalt(SaltLength8)

func Test_aesgcm(t *testing.T) {
	key := DeriveKey([]byte("test-pw"), salt)
	data := []byte("test-data")

	ciphertext, err := EncryptAes(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)
	cleartext, err := DecryptAes(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)

	wrongKey := DeriveKey([]byte("wrong-pw"), salt) // wrong pw, same salt
	cleartext, err = DecryptAes(ciphertext, wrongKey)
	assert.NotNil(t, err, "wrong password should not decrypt")
}

func Test_secretbox(t *testing.T) {
	key := DeriveKey([]byte("test-pw"), salt)
	data := []byte("test-data")

	ciphertext, err := EncryptSecretbox(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)
	cleartext, err := DecryptSecretbox(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)

	wrongKey := DeriveKey([]byte("wrong-pw"), salt) // wrong pw, same salt
	cleartext, err = DecryptSecretbox(ciphertext, wrongKey)
	assert.NotNil(t, err, "wrong password should not decrypt")
}

func Test_XChaCha20Poly1305(t *testing.T) {
	key := DeriveKey([]byte("test-pw"), salt)
	data := []byte("test-data")

	ciphertext, err := EncryptXChaPo(data, key)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)
	cleartext, err := DecryptXChaPo(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)

	wrongKey := DeriveKey([]byte("wrong-pw"), salt) // wrong pw, same salt
	cleartext, err = DecryptXChaPo(ciphertext, wrongKey)
	assert.NotNil(t, err, "wrong password should not decrypt")
}
