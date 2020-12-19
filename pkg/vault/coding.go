package vault

import (
	"encoding/base32"
	"log"
)

// Base32 decoding (only uppercase, no special characters)
func decode(encoded string) []byte {
	decoded, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return decoded
}

// Base32 encoding (only uppercase, no special characters)
func encode(data []byte) string {
	encoded := base32.StdEncoding.EncodeToString(data)
	return encoded
}
