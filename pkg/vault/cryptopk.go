package vault

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/box"
)

type BoxKeyPair struct {
	PrivateKey [32]byte
	PublicKey  [32]byte
}

func (bkp BoxKeyPair) debug() string {
	return "public: " + hex.EncodeToString(bkp.PublicKey[:]) + " private: " + hex.EncodeToString(bkp.PrivateKey[:])
}

func newBoxKeyPair() BoxKeyPair {
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return BoxKeyPair{
		PrivateKey: *privateKey,
		PublicKey:  *publicKey,
	}
}

func newNonce() [24]byte {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func newSecretboxKey() [32]byte {
	var nonce [32]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

type BoxKeys struct {
	OwnPrivateKey  [32]byte
	OtherPublicKey [32]byte
	Nonce          [24]byte
}

// not using 'SealAnonymous' for Box directly, to be able to work with the sender keys and nonce directly
// a comparable approach with ephemeral sender keys is used, however with explicit nonce instead of hashes
// always generates a new "data encryption key", used for NaCl Secretbox, and only encrypts this with NaCl Box
func encryptNaclBox(data []byte, boxKeys BoxKeys) ([]byte, error) {
	var out = make([]byte, 1)
	encrypted := box.Seal(out, data, &boxKeys.Nonce, &boxKeys.OtherPublicKey, &boxKeys.OwnPrivateKey)
	return encrypted[1:], nil
}

func decryptNaclBox(data []byte, boxKeys BoxKeys) ([]byte, error) {
	var out = make([]byte, 1)
	decrypted, ok := box.Open(out, data, &boxKeys.Nonce, &boxKeys.OtherPublicKey, &boxKeys.OwnPrivateKey)
	if !ok {
		err := errors.New("Failed to decrypt box.")
		return nil, err
	}
	return decrypted[1:], nil
}

// combination of Box and Secretbox, starting point for multiple recipients
// currently all output parts are just appended
func encryptNaclBoxSecretBox(data []byte, boxKeys BoxKeys) ([]byte, error) {
	dek := newSecretboxKey()
	iv := newNonce()
	key := &key{
		aesKey: dek[:],
		iv:     iv[:],
	}
	encryptedData, err := encryptSecretbox(data, key)
	if err != nil {
		return nil, err
	}
	encryptedDek, err := encryptNaclBox(dek[:], boxKeys)
	if err != nil {
		return nil, err
	}
	output := encryptedDek
	output = append(output, iv[:]...)
	output = append(output, encryptedData[:]...)
	return output, nil
}

func decryptNaclBoxSecretbox(data []byte, boxKeys BoxKeys) ([]byte, error) {
	keyLength := 32
	ivLength := 24
	encryptedDek := data[:keyLength+box.Overhead]
	iv := data[keyLength+box.Overhead : keyLength+box.Overhead+ivLength]
	encryptedData := data[keyLength+box.Overhead+ivLength:]
	decryptedDek, err := decryptNaclBox(encryptedDek, boxKeys)
	if err != nil {
		return nil, err
	}
	key := &key{
		aesKey: decryptedDek,
		iv:     iv,
	}
	decryptedData, err := decryptSecretbox(encryptedData, key)
	return decryptedData, err
}
