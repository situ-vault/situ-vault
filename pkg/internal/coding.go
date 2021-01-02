package internal

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/nicksnyder/basen"
)

// Base32 decoding (only uppercase, no special characters)
func DecodeBase32(encoded string) ([]byte, error) {
	decoded, err := base32.StdEncoding.DecodeString(encoded)
	return decoded, err
}

// Base32 encoding (only uppercase, no special characters)
func EncodeBase32(data []byte) string {
	encoded := base32.StdEncoding.EncodeToString(data)
	return encoded
}

// Base64 decoding
func DecodeBase64(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	return decoded, err
}

// Base64 encoding
func EncodeBase64(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

// Base64 decoding (URL safe variant)
func DecodeBase64Url(encoded string) ([]byte, error) {
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	return decoded, err
}

// Base64 encoding (URL safe variant)
func EncodeBase64Url(data []byte) string {
	encoded := base64.URLEncoding.EncodeToString(data)
	return encoded
}

// Base62 decoding
func DecodeBase62(encoded string) ([]byte, error) {
	decoded, err := basen.Base62Encoding.DecodeString(encoded)
	return decoded, err
}

// Base62 encoding
func EncodeBase62(data []byte) string {
	encoded := basen.Base62Encoding.EncodeToString(data)
	return encoded
}

// Hex decoding
func DecodeHex(encoded string) ([]byte, error) {
	decoded, err := hex.DecodeString(encoded)
	return decoded, err
}

// Hex encoding
func EncodeHex(data []byte) string {
	encoded := strings.ToUpper(hex.EncodeToString(data))
	return encoded
}
