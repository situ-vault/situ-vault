package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newSalt(t *testing.T) {
	salt := newSalt()
	assert.Len(t, salt, saltLength)
}
