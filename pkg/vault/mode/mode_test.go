package mode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultMode = DefaultConservative
var defaultModeText = "C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32"

var secretboxMode = &Mode{
	Construct: NACL_SECRETBOX,
	Kdf:       ARGON2ID_T1_M65536_C4,
	Salt:      R16B,
	Encoding:  BASE64URL,
}
var secretboxModeText = "C:NACL_SECRETBOX#KDF:ARGON2ID_T1_M65536_C4#SALT:R16B#ENC:BASE64URL"

func Test_ModeText(t *testing.T) {
	text := defaultMode.Text()
	assert.Equal(t, defaultModeText, text)
}

func Test_ModeText2(t *testing.T) {
	text := secretboxMode.Text()
	assert.Equal(t, secretboxModeText, text)
}

func Test_TextMode(t *testing.T) {
	m := NewMode(defaultModeText)
	assert.Equal(t, defaultMode, m)
}

func Test_TextMode2(t *testing.T) {
	m := NewMode(secretboxModeText)
	assert.Equal(t, secretboxMode, m)
}
