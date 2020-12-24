package vault

import (
	"errors"
	"strings"

	"github.com/polarctos/situ-vault/pkg/vault/mode"
)

const (
	prefix    = "SITU_VAULT_V1"
	separator = "##"
)

type Message struct {
	Prefix     string
	Mode       mode.Mode
	Salt       string
	Ciphertext string
}

func buildEnvelope(salt string, cipherText string) string {
	m := mode.DefaultConservative
	return prefix + separator + m.Text() + separator + salt + separator + cipherText
}

func openEnvelope(data string) (salt string, cipherText string, err error) {
	split := strings.Split(data, separator)
	if split[0] != prefix {
		return "", "", errors.New("unknown input")
	}
	m := mode.DefaultConservative
	if mode.NewMode(split[1]).DeepEqual(m) {
		return "", "", errors.New("unknown mode")
	}
	return split[2], split[3], nil
}
