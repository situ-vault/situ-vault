package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"golang.org/x/crypto/nacl/secretbox"
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
	ciphertext := aesGcm.Seal(nil, key.iv, data, nil)
	// tag is contained in the result as the last bytes
	return ciphertext, err
}

func decrypt(data []byte, key *key) ([]byte, error) {
	block, err := aes.NewCipher(key.aesKey)
	if err != nil {
		log.Fatal(err)
	}
	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesGcm.Open(nil, key.iv, data, nil)
	return plaintext, err
}

func encryptSecretbox(data []byte, key *key) ([]byte, error) {
	var secretKey [32]byte
	copy(secretKey[:], key.aesKey)
	var nonce [24]byte
	copy(nonce[:], key.iv)
	copy(nonce[12:], key.iv) // FIXME: golang documented nonce length: 24 bytes, we provide only 12
	// often the result is appended to the nonce, here the nonce comes from the kdf instead
	// however the API kind of expects appending, thus providing 1 byte and slicing afterwards
	var out = make([]byte, 1)
	encrypted := secretbox.Seal(out, data, &nonce, &secretKey)
	return encrypted[1:], nil
}

func decryptSecretbox(data []byte, key *key) ([]byte, error) {
	var secretKey [32]byte
	copy(secretKey[:], key.aesKey)
	var nonce [24]byte
	copy(nonce[:], key.iv)
	copy(nonce[12:], key.iv)
	var out = make([]byte, 1)
	decrypted, ok := secretbox.Open(out, data, &nonce, &secretKey)
	var err error
	if !ok {
		err = errors.New("Failed to decrypt secretbox.")
	}
	return decrypted[1:], err
}
