package vault

import (
	"errors"

	"github.com/polarctos/situ-vault/pkg/internal"
	"github.com/polarctos/situ-vault/pkg/vault/vaultmessage"
	"github.com/polarctos/situ-vault/pkg/vault/vaultmode"
)

func Encrypt(cleartext string, password string, modeText string) (messageText string, err error) {
	mm := vaultmode.NewMode(modeText)

	var salt []byte
	switch mm.Salt {
	case vaultmode.R8B:
		salt = internal.NewSalt(internal.SaltLength8)
	case vaultmode.R16B:
		salt = internal.NewSalt(internal.SaltLength16)
	default:
		return "", errors.New("selected salt variant not implemented")
	}

	pw := []byte(password)
	var key *internal.Key
	switch mm.Kdf {
	case vaultmode.PBKDF2_SHA256_I10K:
		key = internal.DeriveKey(pw, salt)
	case vaultmode.ARGON2ID_T1_M65536_C4:
		key = internal.DeriveKeyArgon2id(pw, salt)
	case vaultmode.SCRYPT_N32768_R8_P1:
		key = internal.DeriveKeyScrypt(pw, salt)
	default:
		return "", errors.New("selected key derivation function not implemented")
	}

	data := []byte(cleartext)
	var encrypted []byte
	switch mm.Construct {
	case vaultmode.AES256_GCM:
		encrypted, err = internal.EncryptAes(data, key)
	case vaultmode.NACL_SECRETBOX:
		encrypted, err = internal.EncryptSecretbox(data, key)
	case vaultmode.XCHACHA20_POLY1305:
		encrypted, err = internal.EncryptXChaPo(data, key)
	default:
		return "", errors.New("selected construct not implemented")
	}

	var encodedSalt string
	var encodedCiphertext string
	switch mm.Encoding {
	case vaultmode.BASE32:
		encodedSalt = internal.EncodeBase32(salt)
		encodedCiphertext = internal.EncodeBase32(encrypted)
	case vaultmode.BASE62:
		encodedSalt = internal.EncodeBase62(salt)
		encodedCiphertext = internal.EncodeBase62(encrypted)
	case vaultmode.BASE64:
		encodedSalt = internal.EncodeBase64(salt)
		encodedCiphertext = internal.EncodeBase64(encrypted)
	case vaultmode.BASE64URL:
		encodedSalt = internal.EncodeBase64Url(salt)
		encodedCiphertext = internal.EncodeBase64Url(encrypted)
	default:
		return "", errors.New("selected encoding not implemented")
	}

	m := vaultmessage.Message{
		Prefix:     vaultmessage.VaultPrefix,
		Mode:       *mm,
		Salt:       encodedSalt,
		Ciphertext: encodedCiphertext,
	}
	messageText = m.Text()
	return messageText, err
}

func Decrypt(messageText string, password string) (cleartext string, modeText string, err error) {
	var message *vaultmessage.Message
	message, err = vaultmessage.NewMessage(messageText)
	if err != nil {
		return "", "", err
	}
	mm := message.Mode

	var decodedSalt []byte
	var decodedCiphertext []byte
	switch mm.Encoding {
	case vaultmode.BASE32:
		decodedSalt, err = internal.DecodeBase32(message.Salt)
		decodedCiphertext, err = internal.DecodeBase32(message.Ciphertext)
	case vaultmode.BASE62:
		decodedSalt, err = internal.DecodeBase62(message.Salt)
		decodedCiphertext, err = internal.DecodeBase62(message.Ciphertext)
	case vaultmode.BASE64:
		decodedSalt, err = internal.DecodeBase64(message.Salt)
		decodedCiphertext, err = internal.DecodeBase64(message.Ciphertext)
	case vaultmode.BASE64URL:
		decodedSalt, err = internal.DecodeBase64Url(message.Salt)
		decodedCiphertext, err = internal.DecodeBase64Url(message.Ciphertext)
	}
	if err != nil {
		return "", "", err
	}

	pw := []byte(password)
	var key *internal.Key
	switch mm.Kdf {
	case vaultmode.PBKDF2_SHA256_I10K:
		key = internal.DeriveKey(pw, decodedSalt)
	case vaultmode.ARGON2ID_T1_M65536_C4:
		key = internal.DeriveKeyArgon2id(pw, decodedSalt)
	case vaultmode.SCRYPT_N32768_R8_P1:
		key = internal.DeriveKeyScrypt(pw, decodedSalt)
	default:
		return "", "", errors.New("selected key derivation function not implemented")
	}

	var decrypted []byte
	switch mm.Construct {
	case vaultmode.AES256_GCM:
		decrypted, err = internal.DecryptAes(decodedCiphertext, key)
	case vaultmode.NACL_SECRETBOX:
		decrypted, err = internal.DecryptSecretbox(decodedCiphertext, key)
	case vaultmode.XCHACHA20_POLY1305:
		decrypted, err = internal.DecryptXChaPo(decodedCiphertext, key)
	default:
		return "", "", errors.New("selected construct not implemented")
	}

	return string(decrypted), mm.Text(), err
}
