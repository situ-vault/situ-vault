package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
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
	copy(nonce[12:], key.iv) // FIXME
	var out = make([]byte, 1)
	decrypted, ok := secretbox.Open(out, data, &nonce, &secretKey)
	var err error
	if !ok {
		err = errors.New("Failed to decrypt secretbox.")
	}
	return decrypted[1:], err
}

// XChaCha20-Poly1305
func encryptXChaPo(data []byte, key *key) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key.aesKey)
	if err != nil {
		panic(err)
	}
	var nonce [24]byte
	copy(nonce[:], key.iv)
	copy(nonce[12:], key.iv) // FIXME: needed nonce length: 24 bytes, we provide only 12
	var out = make([]byte, 1)
	encrypted := aead.Seal(out, nonce[:], data, nil)
	return encrypted[1:], err
}

func decryptXChaPo(data []byte, key *key) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key.aesKey)
	if err != nil {
		panic(err)
	}
	var nonce [24]byte
	copy(nonce[:], key.iv)
	copy(nonce[12:], key.iv) // FIXME
	var out = make([]byte, 1)
	decrypted, err := aead.Open(out, nonce[:], data, nil)
	return decrypted[1:], err
}
