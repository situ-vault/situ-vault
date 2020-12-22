package vault

import (
	"crypto/sha256"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength = 32 // for AES-256
	ivLength  = 12 // common for AES-GCM
)

type key struct {
	aesKey []byte
	iv     []byte
}

// derive key using PBKDF2 with SHA-256
func deriveKey(password []byte, salt []byte) *key {
	// reasonable defaults used by OpenSSL
	var hash = sha256.New
	var iterations = 10000
	kdfResult := pbkdf2.Key(password, salt, iterations, keyLength+ivLength, hash)
	return splitKdfResult(kdfResult)
}

func deriveKeyArgon2id(password []byte, salt []byte) *key {
	// recommended values from the RFC and of the golang documentation
	var time uint32 = 1
	var memory uint32 = 64 * 1024
	var threads uint8 = 4
	kdfResult := argon2.IDKey(password, salt, time, memory, threads, keyLength+ivLength)
	return splitKdfResult(kdfResult)
}

func splitKdfResult(k []byte) *key {
	return &key{
		aesKey: k[:keyLength],
		iv:     k[keyLength : keyLength+ivLength],
	}
}
