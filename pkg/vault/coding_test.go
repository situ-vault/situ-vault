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
}

func Test_base64U(t *testing.T) {
	data := []byte("test-data")
	e := encodeBase64U(data)
	assert.NotEqual(t, e, data)
	d, err := decodeBase64U(e)
	assert.Nil(t, err)
	assert.Equal(t, d, data)
}
