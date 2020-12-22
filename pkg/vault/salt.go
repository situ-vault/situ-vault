package vault

import (
	"crypto/rand"
	"log"
)

const (
	saltLength = 8 // size also used by OpenSSL, but 16 is also common
)

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
