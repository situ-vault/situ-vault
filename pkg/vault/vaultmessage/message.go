package vaultmessage

import (
	"errors"
	"strings"

	"github.com/polarctos/situ-vault/pkg/vault/vaultmode"
)

const (
	VaultPrefix Prefix = "SITU_VAULT_V1"
	separator   string = "##"
)

type Prefix string

func (p Prefix) Text() string {
	return string(p)
}

type Message struct {
	Prefix     Prefix
	Mode       vaultmode.Mode
	Salt       string
	Ciphertext string
}

func New(mode vaultmode.Mode, salt string, ciphertext string) *Message {
	return &Message{
		Prefix:     VaultPrefix,
		Mode:       mode,
		Salt:       salt,
		Ciphertext: ciphertext,
	}
}

func (m Message) Text() string {
	return VaultPrefix.Text() + separator + m.Mode.Text() + separator + m.Salt + separator + m.Ciphertext
}

func NewMessage(data string) (*Message, error) {
	clean := strings.TrimSpace(data)
	split := strings.Split(clean, separator)
	if len(split) != 4 {
		return nil, errors.New("unknown input length")
	}
	if split[0] != VaultPrefix.Text() {
		return nil, errors.New("unknown input prefix")
	}
	mode, err := vaultmode.NewMode(split[1])
	if err != nil {
		return nil, err
	}

	m := Message{
		Prefix:     Prefix(split[0]),
		Mode:       *mode,
		Salt:       split[2],
		Ciphertext: split[3],
	}
	return &m, nil
}
