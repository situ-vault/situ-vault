package vaultmode

import (
	"errors"
	"reflect"
)

type Construct string

type constructs struct {
	Aes256gcm         Construct
	NaclSecretbox     Construct
	XChaCha20poly1305 Construct
}

var Constructs = constructs{
	Aes256gcm:         Aes256gcm,
	NaclSecretbox:     NaclSecretbox,
	XChaCha20poly1305: XChaCha20poly1305,
}

const (
	Aes256gcm         Construct = "AES256_GCM"         // AEAD; Standard: Nonce: 12 byte, Tag: 16 byte (year 2000/2007)
	NaclSecretbox     Construct = "NACL_SECRETBOX"     // AE; XSalsa20-Poly1305; Standard: Nonce: 24 byte, Tag: 16 byte (year 2008)
	XChaCha20poly1305 Construct = "XCHACHA20_POLY1305" // AEAD; XChaCha20-Poly1305; Standard: Nonce: 24 byte, Tag: 16 byte (year 2008/2018)
	NaclBox           Construct = "NACL_BOX"           // PKE; Curve25519-XSalsa20-Poly1305; Standard: PublicKey: 32 byte, SecretKey: 24 byte, Nonce: 24 byte, Tag: 16 byte
	NaclBoxSecretbox  Construct = "NACL_BOX_SECRETBOX" // AE & PKE; NACL_SECRETBOX for the actual data, and NACL_BOX for a data encryption key
)

func ParseConstruct(s string) (Construct, error) {
	for _, value := range allConstructValues {
		if s == value {
			return Construct(s), nil
		}
	}
	return "", errors.New("Invalid value: " + s)
}

var allConstructValues = allValues(reflect.ValueOf(Constructs))

func (c constructs) AllValues() []string {
	return allConstructValues
}
