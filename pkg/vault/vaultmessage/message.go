package vaultmessage

import (
	"errors"
	"strings"

	"github.com/situ-vault/situ-vault/pkg/vault/vaultmode"
)

const (
	VaultPrefix Prefix = "SITU_VAULT_V1"
	separator   string = "##"
	end         string = "END"
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
	End        string
}

func New(mode vaultmode.Mode, salt string, ciphertext string) *Message {
	return &Message{
		Prefix:     VaultPrefix,
		Mode:       mode,
		Salt:       salt,
		Ciphertext: ciphertext,
		End:        end,
	}
}

func (m Message) Text() string {
	return VaultPrefix.Text() + separator + m.Mode.Text() + separator + m.Salt + separator + m.Ciphertext + separator + m.End
}

func NewMessage(data string) (*Message, error) {
	clean := strings.TrimSpace(data)
	split := strings.Split(clean, separator)
	if len(split) != 5 {
		return nil, errors.New("unknown input length")
	}
	prefix := split[0]
	if prefix != VaultPrefix.Text() {
		return nil, errors.New("unknown input prefix")
	}
	mode, err := vaultmode.NewMode(split[1])
	if err != nil {
		return nil, err
	}
	salt := split[2]
	if len(salt) < 1 {
		return nil, errors.New("missing salt")
	}
	ciphertext := split[3]
	if len(ciphertext) < 1 {
		return nil, errors.New("missing ciphertext")
	}
	z := split[4]
	if z != end {
		return nil, errors.New("missing end")
	}

	m := Message{
		Prefix:     Prefix(prefix),
		Mode:       *mode,
		Salt:       salt,
		Ciphertext: ciphertext,
		End:        z,
	}
	return &m, nil
}
