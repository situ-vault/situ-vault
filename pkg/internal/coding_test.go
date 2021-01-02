package internal

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_base32(t *testing.T) {
	data := []byte("test-data")
	e := EncodeBase32(data)
	assert.NotEqual(t, e, data)
	d, err := DecodeBase32(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
	assert.Equal(t, "ORSXG5BNMRQXIYI=", e)
}

func Test_base64(t *testing.T) {
	data := []byte("test-data")
	e := EncodeBase64(data)
	assert.NotEqual(t, e, data)
	d, err := DecodeBase64(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
	assert.Equal(t, "dGVzdC1kYXRh", e)
}

func Test_base64Url(t *testing.T) {
	data := []byte("test-data")
	e := EncodeBase64Url(data)
	assert.NotEqual(t, e, data)
	d, err := DecodeBase64Url(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
	assert.Equal(t, "dGVzdC1kYXRh", e)
}

func Test_base62(t *testing.T) {
	data := []byte("test-data")
	e := EncodeBase62(data)
	assert.NotEqual(t, e, data)
	d, err := DecodeBase62(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
	assert.Equal(t, "fGF8D3pR6tsH", e)
}

func Test_hex(t *testing.T) {
	data := []byte("test-data")
	e := EncodeHex(data)
	assert.NotEqual(t, e, data)
	assert.Regexp(t, regexp.MustCompile(`^([0-9A-F]{2})+$`), e, "only uppercase intended")
	d, err := DecodeHex(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
	d, err = DecodeHex(strings.ToLower(e))
	assert.Nil(t, err, "also allow decoding of lowercase")
	assert.Equal(t, d, data)
	assert.Equal(t, "746573742D64617461", e)
}
