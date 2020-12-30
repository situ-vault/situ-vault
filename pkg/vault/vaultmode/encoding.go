package vaultmode

import (
	"errors"
	"reflect"
)

type Encoding string

const (
	None      Encoding = "NONE"      // No encoding, just bytes
	Base32    Encoding = "BASE32"    // Base32
	Base62    Encoding = "BASE62"    // Base62 (is Base64 without the 2 special characters)
	Base64    Encoding = "BASE64"    // Base64
	Base64url Encoding = "BASE64URL" // Base64 (URL safe variant)
)

type encodings struct {
	None      Encoding
	Base32    Encoding
	Base62    Encoding
	Base64    Encoding
	Base64url Encoding
}

var Encodings = encodings{
	None:      None,
	Base32:    Base32,
	Base62:    Base62,
	Base64:    Base64,
	Base64url: Base64url,
}

func ParseEncoding(s string) (Encoding, error) {
	for _, value := range allEncodingValues {
		if s == value {
			return Encoding(s), nil
		}
	}
	return "", errors.New("Invalid value: " + s)
}

var allEncodingValues = allValues(reflect.ValueOf(Encodings))

func (e encodings) AllValues() []string {
	return allEncodingValues
}
