package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NaclBox(t *testing.T) {
	testPk(t, encryptNaclBox, decryptNaclBox)
}

func Test_NaclBoxSecretbox(t *testing.T) {
	testPk(t, encryptNaclBoxSecretBox, decryptNaclBoxSecretbox)
}

func testPk(
	t *testing.T,
	enc func(data []byte, boxKeys BoxKeys) ([]byte, error),
	dec func(data []byte, boxKeys BoxKeys) ([]byte, error)) {
	data := []byte("test-data")
	recipientKeys := newBoxKeyPair()
	ephemeralSenderKeys := newBoxKeyPair()

	boxKeysEncrypt := BoxKeys{
		OwnPrivateKey:  ephemeralSenderKeys.PrivateKey,
		OtherPublicKey: recipientKeys.PublicKey,
		Nonce:          newNonce(),
	}
	ciphertext, err := enc(data, boxKeysEncrypt)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)

	moreKeys := newBoxKeyPair()
	boxKeysMore := BoxKeys{
		OwnPrivateKey:  moreKeys.PrivateKey,           // different recipient
		OtherPublicKey: ephemeralSenderKeys.PublicKey, // same sender
		Nonce:          boxKeysEncrypt.Nonce,          // same nonce
	}
	_, err = dec(ciphertext, boxKeysMore)
	assert.NotNil(t, err, "wrong keys should not decrypt")

	boxKeysDecrypt := BoxKeys{
		OwnPrivateKey:  recipientKeys.PrivateKey,      // correct recipient
		OtherPublicKey: ephemeralSenderKeys.PublicKey, // same sender
		Nonce:          boxKeysEncrypt.Nonce,          // same nonce
	}
	cleartext, err := dec(ciphertext, boxKeysDecrypt)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)
}
