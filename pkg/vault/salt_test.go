package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newSalt(t *testing.T) {
	salt := newSalt()
	assert.Len(t, salt, saltLength)
}
