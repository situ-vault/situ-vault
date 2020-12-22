package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testPassword = []byte("test-password")
var testSalt = []byte("test-salt")

var expectedKey = key{
	aesKey: []byte{0x20, 0xb3, 0xc, 0x92, 0x34, 0xc1, 0x97, 0x8f, 0x33, 0xc, 0xe4, 0x9d, 0x47, 0x23, 0x85, 0x41, 0xe2, 0x2f, 0x5a, 0xcc, 0x77, 0x6f, 0xd5, 0x8f, 0x66, 0x4a, 0x55, 0xf4, 0x70, 0xa4, 0xa5, 0x52},
	iv:     []byte{0xd4, 0x69, 0xa8, 0xde, 0x62, 0xb9, 0x54, 0xab, 0xb, 0x6b, 0x75, 0x9c},
}
var expectedKeyArgon2id = key{
	aesKey: []byte{0x4f, 0xa4, 0x77, 0x60, 0x5a, 0x8, 0x61, 0x90, 0x7e, 0x61, 0x1f, 0x2b, 0x4d, 0x8f, 0x88, 0x21, 0x82, 0x87, 0x6b, 0x55, 0xbd, 0x55, 0xf5, 0x4a, 0xf5, 0x36, 0xc3, 0xec, 0xa9, 0x9f, 0xa5, 0xeb},
	iv:     []byte{0x4e, 0xdc, 0x35, 0x68, 0x10, 0x5f, 0x90, 0xec, 0xa1, 0x89, 0x74, 0x86},
}

func Test_deriveKeyPbkdf2(t *testing.T) {
	key := deriveKey(testPassword, testSalt)

	assert.Len(t, key.aesKey, keyLength)
	assert.Len(t, key.iv, ivLength)

	assert.EqualValues(t, expectedKey.aesKey, key.aesKey)
	assert.EqualValues(t, expectedKey.iv, key.iv)
}

func Test_deriveKeyArgon2(t *testing.T) {
	key := deriveKeyArgon2id(testPassword, testSalt)

	assert.Len(t, key.aesKey, keyLength)
	assert.Len(t, key.iv, ivLength)

	assert.EqualValues(t, expectedKeyArgon2id.aesKey, key.aesKey)
	assert.EqualValues(t, expectedKeyArgon2id.iv, key.iv)
}
