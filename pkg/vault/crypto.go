package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
)

func encrypt(data []byte, key *key) ([]byte, error) {
	block, err := aes.NewCipher(key.aesKey)
	if err != nil {
		log.Fatal(err)
	}
	// standard nonce length: 12 bytes
	// standard tag length: 16 bytes
	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	ciphertext := aesGcm.Seal(nil, key.iv, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, key *key) ([]byte, error) {
	block, err := aes.NewCipher(key.aesKey)
	if err != nil {
		log.Fatal(err)
	}
	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	plaintext, err := aesGcm.Open(nil, key.iv, data, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext, nil
}
