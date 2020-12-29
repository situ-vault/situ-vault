package vaultmode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultMode = Defaults().Conservative
var defaultModeText = "C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32"

var secretboxMode = Defaults().Secretbox
var secretboxModeText = "C:NACL_SECRETBOX#KDF:SCRYPT_N32768_R8_P1#SALT:R16B#ENC:BASE64"

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
	assert.Equal(t, &defaultMode, m)
}

func Test_TextMode2(t *testing.T) {
	m := NewMode(secretboxModeText)
	assert.Equal(t, &secretboxMode, m)
}
