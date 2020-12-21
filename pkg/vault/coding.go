package vault

import (
	"encoding/base32"
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
