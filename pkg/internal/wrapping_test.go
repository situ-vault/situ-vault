package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wrap1(t *testing.T) {
	wrapped := Wrap("ab", 1)
	assert.Equal(t, newlineChars+"a"+newlineChars+"b", wrapped)
}

func Test_wrap2(t *testing.T) {
	wrapped := Wrap("abcdefg", 2)
	assert.Equal(t, newlineChars+"ab"+newlineChars+"cd"+newlineChars+"ef"+newlineChars+"g", wrapped)
}

func Test_unwrap(t *testing.T) {
	unwrapped := Unwrap("\na\nb\r\nc")
	assert.Equal(t, "abc", unwrapped)
}
