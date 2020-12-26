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

type BoxKeys struct {
	OwnPrivateKey  [32]byte
	OtherPublicKey [32]byte
	Nonce          [24]byte
}

// not using 'SealAnonymous' directly, to be able to work with the sender keys and nonce directly
// a comparable approach with ephemeral sender keys is used, however with explicit nonce instead of hashes
func encryptNaclBox(data []byte, boxKeys BoxKeys) ([]byte, error) {
	var out = make([]byte, 1)
	encrypted := box.Seal(out, data, &boxKeys.Nonce, &boxKeys.OtherPublicKey, &boxKeys.OwnPrivateKey)
	return encrypted[1:], nil
}

func decryptNaclBox(data []byte, boxKeys BoxKeys) ([]byte, error) {
	var out = make([]byte, 1)
	decrypted, ok := box.Open(out, data, &boxKeys.Nonce, &boxKeys.OtherPublicKey, &boxKeys.OwnPrivateKey)
	var err error
	if !ok {
		err = errors.New("Failed to decrypt box.")
		return nil, err
	}
	return decrypted[1:], nil
}
