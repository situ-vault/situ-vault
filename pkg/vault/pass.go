package vault

import (
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
	"log"
)

const (
	keyLength  = 32    // for AES-256
	ivLength   = 12    // common for AES-GCM
	iterations = 10000 // reasonable default of OpenSSL
	saltLength = 8     // size also used by OpenSSL
)

type key struct {
	aesKey []byte
	iv     []byte
}

// derive key using PBKDF2 with SHA-256
func deriveKey(password []byte, salt []byte) *key {
	kdfResult := pbkdf2.Key(password, salt, iterations, keyLength+ivLength, sha256.New)
	return &key{
		aesKey: kdfResult[:keyLength],
		iv:     kdfResult[keyLength : keyLength+ivLength],
	}
}

func newSalt() []byte {
	salt, err := randomBytes(saltLength)
	if err != nil {
		log.Fatal(err)
	}
	return salt
}

func randomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	return b, err
}
