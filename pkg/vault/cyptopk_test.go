package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NaclBox(t *testing.T) {
	data := []byte("test-data")
	recipientKeys := newBoxKeyPair()
	ephemeralSenderKeys := newBoxKeyPair()

	boxKeysEncrypt := BoxKeys{
		OwnPrivateKey:  ephemeralSenderKeys.PrivateKey,
		OtherPublicKey: recipientKeys.PublicKey,
		Nonce:          newNonce(),
	}
	ciphertext, err := encryptNaclBox(data, boxKeysEncrypt)
	assert.Nil(t, err)
	assert.NotContains(t, ciphertext, data)

	moreKeys := newBoxKeyPair()
	boxKeysMore := BoxKeys{
		OwnPrivateKey:  moreKeys.PrivateKey,           // different recipient
		OtherPublicKey: ephemeralSenderKeys.PublicKey, // same sender
		Nonce:          boxKeysEncrypt.Nonce,          // same nonce
	}
	_, err = decryptNaclBox(ciphertext, boxKeysMore)
	assert.NotNil(t, err, "wrong keys should not decrypt")

	boxKeysDecrypt := BoxKeys{
		OwnPrivateKey:  recipientKeys.PrivateKey,      // correct recipient
		OtherPublicKey: ephemeralSenderKeys.PublicKey, // same sender
		Nonce:          boxKeysEncrypt.Nonce,          // same nonce
	}
	cleartext, err := decryptNaclBox(ciphertext, boxKeysDecrypt)
	assert.Nil(t, err)
	assert.Equal(t, cleartext, data)
}
