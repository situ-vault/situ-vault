package internal

import (
	"crypto/rand"
)

const (
	SaltLength8  = 8 // size also used by OpenSSL, but 16 is also common
	SaltLength16 = 16
)

func NewSalt(length int) []byte {
	salt, err := randomBytes(length)
	if err != nil {
		panic(err)
	}
	return salt
}

func randomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	return b, err
}
