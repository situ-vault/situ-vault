package vault

import (
	"encoding/base32"
	"encoding/base64"

	"github.com/nicksnyder/basen"
)

// Base32 decoding (only uppercase, no special characters)
func decode(encoded string) ([]byte, error) {
	decoded, err := base32.StdEncoding.DecodeString(encoded)
	return decoded, err
}

// Base32 encoding (only uppercase, no special characters)
func encode(data []byte) string {
	encoded := base32.StdEncoding.EncodeToString(data)
	return encoded
}

// Base64 decoding (URL safe variant)
func decodeBase64U(encoded string) ([]byte, error) {
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	return decoded, err
}

// Base64 encoding (URL safe variant)
func encodeBase64U(data []byte) string {
	encoded := base64.URLEncoding.EncodeToString(data)
	return encoded
}

// Base62 decoding
func decodeBase62(encoded string) ([]byte, error) {
	decoded, err := basen.Base62Encoding.DecodeString(encoded)
	return decoded, err
}

// Base62 encoding
func encodeBase62(data []byte) string {
	encoded := basen.Base62Encoding.EncodeToString(data)
	return encoded
}
