package internal

import (
	"crypto/sha256"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

const (
	keyLength = 32 // for AES-256
	ivLength  = 24 // 12 common for AES-GCM but 24 for XSalsa20 and XChaCha
)

type Key struct {
	key []byte
	iv  []byte
}

// derive key using PBKDF2 with SHA-256
func DeriveKeyPbkdf(password []byte, salt []byte) *Key {
	// reasonable defaults used by OpenSSL
	var hash = sha256.New
	var iterations = 10000
	kdfResult := pbkdf2.Key(password, salt, iterations, keyLength+ivLength, hash)
	return splitKdfResult(kdfResult)
}

func DeriveKeyArgon2id(password []byte, salt []byte) *Key {
	// recommended values from the RFC and of the golang documentation
	var time uint32 = 1
	var memory uint32 = 64 * 1024
	var threads uint8 = 4
	kdfResult := argon2.IDKey(password, salt, time, memory, threads, keyLength+ivLength)
	return splitKdfResult(kdfResult)
}

func DeriveKeyScrypt(password []byte, salt []byte) *Key {
	// recommended values from the golang documentation
	var N = 1 << 15 // N is a CPU/memory cost parameter, power of two greater than 1. here 32768
	var r = 8       // r and p must satisfy r * p < 2³⁰
	var p = 1
	kdfResult, err := scrypt.Key(password, salt, N, r, p, keyLength+ivLength)
	if err != nil {
		panic(err)
	}
	return splitKdfResult(kdfResult)
}

// derive key using HKDF with SHA-256
func DeriveKeyHkdf(password []byte, salt []byte) *Key {
	var hash = sha256.New
	var info []byte = nil // optional, not used
	hkdf := hkdf.New(hash, password, salt, info)
	length := keyLength + ivLength
	kdfResult := make([]byte, length)
	read, err := hkdf.Read(kdfResult)
	if err != nil || read != length {
		panic(err)
	}
	return splitKdfResult(kdfResult)
}

func splitKdfResult(k []byte) *Key {
	return &Key{
		key: k[:keyLength],
		iv:  k[keyLength : keyLength+ivLength],
	}
}
