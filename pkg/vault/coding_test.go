package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_base32(t *testing.T) {
	data := []byte("test-data")
	e := encode(data)
	assert.NotEqual(t, e, data)
	d, err := decode(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
	assert.Equal(t, "ORSXG5BNMRQXIYI=", e)
}

func Test_base64U(t *testing.T) {
	data := []byte("test-data")
	e := encodeBase64U(data)
	assert.NotEqual(t, e, data)
	d, err := decodeBase64U(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
	assert.Equal(t, "dGVzdC1kYXRh", e)
}

func Test_base62(t *testing.T) {
	data := []byte("test-data")
	e := encodeBase62(data)
	assert.NotEqual(t, e, data)
	d, err := decodeBase62(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
	assert.Equal(t, "fGF8D3pR6tsH", e)
}
