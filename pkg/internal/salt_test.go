package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newSalt(t *testing.T) {
	salt := NewSalt(SaltLength8)
	assert.Len(t, salt, SaltLength8)

	salt2 := NewSalt(SaltLength8)
	assert.NotEqual(t, salt, salt2, "always a new salt")
}
