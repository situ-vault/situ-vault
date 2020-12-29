package vault

import (
	"errors"

	"github.com/polarctos/situ-vault/pkg/internal"
	"github.com/polarctos/situ-vault/pkg/vault/vaultmessage"
	"github.com/polarctos/situ-vault/pkg/vault/vaultmode"
)

func Encrypt(cleartext string, password string, modeText string) (messageText string, err error) {
	mm, err := vaultmode.NewMode(modeText)
	if err != nil {
		return "", err
	}

	var salt []byte
	switch mm.Salt {
	case vaultmode.Salts.R8b:
		salt = internal.NewSalt(internal.SaltLength8)
	case vaultmode.Salts.R16b:
		salt = internal.NewSalt(internal.SaltLength16)
	default:
		return "", errors.New("selected salt variant not implemented")
	}

	pw := []byte(password)
	var key *internal.Key
	switch mm.Kdf {
	case vaultmode.KeyDerivationFunctions.Pbkdf2_sha256_i10k:
		key = internal.DeriveKey(pw, salt)
	case vaultmode.KeyDerivationFunctions.Argon2id_t1_m65536_c4:
		key = internal.DeriveKeyArgon2id(pw, salt)
	case vaultmode.KeyDerivationFunctions.Scrypt_n32768_r8_p1:
		key = internal.DeriveKeyScrypt(pw, salt)
	default:
		return "", errors.New("selected key derivation function not implemented")
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

	var encodedSalt string
	var encodedCiphertext string
	switch mm.Encoding {
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

	m := vaultmessage.New(*mm, encodedSalt, encodedCiphertext)
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
	case vaultmode.Encodings.Base32:
		decodedSalt, err = internal.DecodeBase32(message.Salt)
		decodedCiphertext, err = internal.DecodeBase32(message.Ciphertext)
	case vaultmode.Encodings.Base62:
		decodedSalt, err = internal.DecodeBase62(message.Salt)
		decodedCiphertext, err = internal.DecodeBase62(message.Ciphertext)
	case vaultmode.Encodings.Base64:
		decodedSalt, err = internal.DecodeBase64(message.Salt)
		decodedCiphertext, err = internal.DecodeBase64(message.Ciphertext)
	case vaultmode.Encodings.Base64url:
		decodedSalt, err = internal.DecodeBase64Url(message.Salt)
		decodedCiphertext, err = internal.DecodeBase64Url(message.Ciphertext)
	}
	if err != nil {
		return "", "", err
	}

	pw := []byte(password)
	var key *internal.Key
	switch mm.Kdf {
	case vaultmode.KeyDerivationFunctions.Pbkdf2_sha256_i10k:
		key = internal.DeriveKey(pw, decodedSalt)
	case vaultmode.KeyDerivationFunctions.Argon2id_t1_m65536_c4:
		key = internal.DeriveKeyArgon2id(pw, decodedSalt)
	case vaultmode.KeyDerivationFunctions.Scrypt_n32768_r8_p1:
		key = internal.DeriveKeyScrypt(pw, decodedSalt)
	default:
		return "", "", errors.New("selected key derivation function not implemented")
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
