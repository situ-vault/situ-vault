package vault

import (
	"errors"
	"strings"

	"github.com/polarctos/situ-vault/pkg/vault/mode"
)

const (
	prefix    Prefix = "SITU_VAULT_V1"
	separator string = "##"
)

type Prefix string

func (p Prefix) Text() string {
	return string(p)
}

type Message struct {
	Prefix     Prefix
	Mode       mode.Mode
	Salt       string
	Ciphertext string
}

func (m Message) Text() string {
	return prefix.Text() + separator + m.Mode.Text() + separator + m.Salt + separator + m.Ciphertext
}

func NewMessage(data string) (*Message, error) {
	split := strings.Split(data, separator)
	if len(split) != 4 {
		return nil, errors.New("unknown input length")
	}
	if split[0] != prefix.Text() {
		return nil, errors.New("unknown input prefix")
	}
	m := Message{
		Prefix:     Prefix(split[0]),
		Mode:       *mode.NewMode(split[1]),
		Salt:       split[2],
		Ciphertext: split[3],
	}
	return &m, nil
}
