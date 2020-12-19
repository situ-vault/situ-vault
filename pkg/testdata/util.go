package testdata

import (
	"crypto/rand"
	"encoding/ascii85"
	"encoding/base64"
)

// Test data utilities

// just for convenience ascii85 formatted password, results in length*5 characters
func RandomPassword(length int) string {
	b, _ := RandomBytes(length * 4)
	dst := make([]byte, length*5)
	ascii85.Encode(dst, b)
	return string(dst)
}

// just for convenience base64 formatted data, length defines the bytes before encoding
func RandomDataBase64(length int) string {
	b, _ := RandomBytes(length)
	return base64.StdEncoding.EncodeToString(b)
}

// array of random bytes, of the requested length
func RandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	return b, err
}
