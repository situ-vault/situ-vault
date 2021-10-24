package vault

import (
	"errors"

	"github.com/situ-vault/situ-vault/pkg/internal"
	"github.com/situ-vault/situ-vault/pkg/vault/vaultmessage"
	"github.com/situ-vault/situ-vault/pkg/vault/vaultmode"
)

func Encrypt(cleartext string, password string, modeText string) (messageText string, err error) {
	cleartext, password, err = maybeFromFiles(cleartext, password)
	if err != nil {
		return "", err
	}

	mm, err := vaultmode.NewMode(modeText)
	if err != nil {
		return "", err
	}

	salt, err := newSalt(mm.Salt)
	if err != nil {
		return "", err
	}

	key, err := deriveKey(mm.Kdf, []byte(password), salt)
	if err != nil {
		return "", err
	}

	data := []byte(cleartext)
	var encrypted []byte
	switch mm.Construct {
	case vaultmode.Constructs.Aes256gcm:
		encrypted, err = internal.EncryptAes(data, key)
	case vaultmode.Constructs.NaclSecretbox:
		encrypted, err = internal.EncryptSecretbox(data, key)
	case vaultmode.Constructs.XChaCha20poly1305:
		encrypted, err = internal.EncryptXChaPo(data, key)
	default:
		return "", errors.New("selected construct not implemented")
	}
	if err != nil {
		return "", err
	}

	var encodedSalt string
	var encodedCiphertext string
	switch mm.Encoding {
	case vaultmode.Encodings.Hex:
		encodedSalt = internal.EncodeHex(salt)
		encodedCiphertext = internal.EncodeHex(encrypted)
	case vaultmode.Encodings.Base32:
		encodedSalt = internal.EncodeBase32(salt)
		encodedCiphertext = internal.EncodeBase32(encrypted)
	case vaultmode.Encodings.Base62:
		encodedSalt = internal.EncodeBase62(salt)
		encodedCiphertext = internal.EncodeBase62(encrypted)
	case vaultmode.Encodings.Base64:
		encodedSalt = internal.EncodeBase64(salt)
		encodedCiphertext = internal.EncodeBase64(encrypted)
	case vaultmode.Encodings.Base64url:
		encodedSalt = internal.EncodeBase64Url(salt)
		encodedCiphertext = internal.EncodeBase64Url(encrypted)
	default:
		return "", errors.New("selected encoding not implemented")
	}

	wrappedCiphertext, err := wrapText(mm.Linebreak, encodedCiphertext)
	if err != nil {
		return "", err
	}

	m := vaultmessage.New(*mm, encodedSalt, wrappedCiphertext)
	messageText = m.Text()
	return messageText, err
}

func Decrypt(messageText string, password string) (cleartext string, modeText string, err error) {
	messageText, password, err = maybeFromFiles(messageText, password)
	if err != nil {
		return "", "", err
	}

	var message *vaultmessage.Message
	message, err = vaultmessage.NewMessage(messageText)
	if err != nil {
		return "", "", err
	}
	mm := message.Mode

	var unwrappedCiphertext string
	if mm.Linebreak != vaultmode.Linebreaks.No {
		unwrappedCiphertext = internal.Unwrap(message.Ciphertext)
	} else {
		unwrappedCiphertext = message.Ciphertext
	}

	var decodedSalt []byte
	var decodedCiphertext []byte
	var errC error
	switch mm.Encoding {
	case vaultmode.Encodings.Hex:
		decodedSalt, err = internal.DecodeHex(message.Salt)
		decodedCiphertext, errC = internal.DecodeHex(unwrappedCiphertext)
	case vaultmode.Encodings.Base32:
		decodedSalt, err = internal.DecodeBase32(message.Salt)
		decodedCiphertext, errC = internal.DecodeBase32(unwrappedCiphertext)
	case vaultmode.Encodings.Base62:
		decodedSalt, err = internal.DecodeBase62(message.Salt)
		decodedCiphertext, errC = internal.DecodeBase62(unwrappedCiphertext)
	case vaultmode.Encodings.Base64:
		decodedSalt, err = internal.DecodeBase64(message.Salt)
		decodedCiphertext, errC = internal.DecodeBase64(unwrappedCiphertext)
	case vaultmode.Encodings.Base64url:
		decodedSalt, err = internal.DecodeBase64Url(message.Salt)
		decodedCiphertext, errC = internal.DecodeBase64Url(unwrappedCiphertext)
	}
	if err != nil {
		return "", "", err
	}
	if errC != nil {
		return "", "", err
	}

	key, err := deriveKey(mm.Kdf, []byte(password), decodedSalt)
	if err != nil {
		return "", "", err
	}

	var decrypted []byte
	switch mm.Construct {
	case vaultmode.Constructs.Aes256gcm:
		decrypted, err = internal.DecryptAes(decodedCiphertext, key)
	case vaultmode.Constructs.NaclSecretbox:
		decrypted, err = internal.DecryptSecretbox(decodedCiphertext, key)
	case vaultmode.Constructs.XChaCha20poly1305:
		decrypted, err = internal.DecryptXChaPo(decodedCiphertext, key)
	default:
		return "", "", errors.New("selected construct not implemented")
	}

	return string(decrypted), mm.Text(), err
}

func newSalt(s vaultmode.Salt) ([]byte, error) {
	var salt []byte
	switch s {
	case vaultmode.Salts.R8b:
		salt = internal.NewSalt(internal.SaltLength8)
	case vaultmode.Salts.R16b:
		salt = internal.NewSalt(internal.SaltLength16)
	case vaultmode.Salts.R24b:
		salt = internal.NewSalt(internal.SaltLength24)
	case vaultmode.Salts.R32b:
		salt = internal.NewSalt(internal.SaltLength32)
	default:
		return nil, errors.New("selected salt variant not implemented")
	}
	return salt, nil
}

func deriveKey(kdf vaultmode.KeyDerivationFunction, pw []byte, salt []byte) (*internal.Key, error) {
	var key *internal.Key
	switch kdf {
	case vaultmode.KeyDerivationFunctions.Pbkdf2_sha256_i10k:
		key = internal.DeriveKeyPbkdf(pw, salt)
	case vaultmode.KeyDerivationFunctions.Argon2id_t1_m65536_c4:
		key = internal.DeriveKeyArgon2id(pw, salt)
	case vaultmode.KeyDerivationFunctions.Scrypt_n32768_r8_p1:
		key = internal.DeriveKeyScrypt(pw, salt)
	case vaultmode.KeyDerivationFunctions.Hkdf_sha256_noinfo:
		key = internal.DeriveKeyHkdf(pw, salt)
	default:
		return nil, errors.New("selected key derivation function not implemented")
	}
	return key, nil
}

func wrapText(linebreak vaultmode.Linebreak, s string) (string, error) {
	switch linebreak {
	case vaultmode.Linebreaks.No:
		return s, nil
	case vaultmode.Linebreaks.Ch80:
		return internal.Wrap(s, 80), nil
	case vaultmode.Linebreaks.Ch100:
		return internal.Wrap(s, 100), nil
	case vaultmode.Linebreaks.Ch120:
		return internal.Wrap(s, 120), nil
	default:
		return "", errors.New("selected line breaking not implemented")
	}
}
