package testdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_randomPassword(t *testing.T) {
	length := 10
	password := RandomPassword(length)
	assert.Len(t, password, length*5)
}

func Test_randomDataBase64(t *testing.T) {
	data := RandomDataBase64(400)
	assert.Len(t, data, 400*1.34)
}
