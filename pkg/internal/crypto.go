package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"log"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/nacl/secretbox"
)

// AES-256-GCM (~ AES-CTR-GMAC)
func EncryptAes(data []byte, key *Key) ([]byte, error) {
	block, err := aes.NewCipher(key.key)
	if err != nil {
		log.Fatal(err)
	}
	// standard nonce length: 12 bytes
	var nonce [12]byte
	copy(nonce[:], key.iv)
	// standard tag length: 16 bytes
	aesGcm, err := cipher.NewGCM(block)
	encrypted := aesGcm.Seal(nil, nonce[:], data, nil)
	// tag is contained in the result as the last bytes
	return encrypted, err
}

func DecryptAes(data []byte, key *Key) ([]byte, error) {
	block, err := aes.NewCipher(key.key)
	if err != nil {
		log.Fatal(err)
	}
	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	var nonce [12]byte
	copy(nonce[:], key.iv)
	decrypted, err := aesGcm.Open(nil, nonce[:], data, nil)
	return decrypted, err
}

// XSalsa20-Poly1305
func EncryptSecretbox(data []byte, key *Key) ([]byte, error) {
	var secretKey [32]byte
	copy(secretKey[:], key.key)
	var nonce [24]byte
	copy(nonce[:], key.iv)
	// often the result is appended to the nonce, here the nonce comes from the kdf instead
	// however the API kind of expects appending, thus providing 1 byte and slicing afterwards
	var out = make([]byte, 1)
	encrypted := secretbox.Seal(out, data, &nonce, &secretKey)
	return encrypted[1:], nil
}

func DecryptSecretbox(data []byte, key *Key) ([]byte, error) {
	var secretKey [32]byte
	copy(secretKey[:], key.key)
	var nonce [24]byte
	copy(nonce[:], key.iv)
	var out = make([]byte, 1)
	decrypted, ok := secretbox.Open(out, data, &nonce, &secretKey)
	var err error
	if !ok {
		err = errors.New("Failed to decrypt secretbox.")
		return nil, err
	}
	return decrypted[1:], nil
}

// XChaCha20-Poly1305
func EncryptXChaPo(data []byte, key *Key) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key.key)
	if err != nil {
		panic(err)
	}
	var nonce [24]byte
	copy(nonce[:], key.iv)
	var out = make([]byte, 1)
	encrypted := aead.Seal(out, nonce[:], data, nil)
	return encrypted[1:], nil
}

func DecryptXChaPo(data []byte, key *Key) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key.key)
	if err != nil {
		panic(err)
	}
	var nonce [24]byte
	copy(nonce[:], key.iv)
	var out = make([]byte, 1)
	decrypted, err := aead.Open(out, nonce[:], data, nil)
	if err != nil {
		return nil, err
	}
	return decrypted[1:], nil
}
