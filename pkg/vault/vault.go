package vault

import (
	"errors"

	"github.com/polarctos/situ-vault/pkg/vault/mode"
)

func Encrypt(cleartext string, password string, modeText string) (messageText string, err error) {
	mm := mode.NewMode(modeText)

	var salt []byte
	switch mm.Salt {
	case mode.R8B:
		salt = newSalt(SaltLength8)
	case mode.R16B:
		salt = newSalt(SaltLength16)
	default:
		return "", errors.New("selected salt variant not implemented")
	}

	pw := []byte(password)
	var key *key
	switch mm.Kdf {
	case mode.PBKDF2_SHA256_I10K:
		key = deriveKey(pw, salt)
	case mode.ARGON2ID_T1_M65536_C4:
		key = deriveKeyArgon2id(pw, salt)
	case mode.SCRYPT_N32768_R8_P1:
		key = deriveKeyScrypt(pw, salt)
	default:
		return "", errors.New("selected key derivation function not implemented")
	}

	data := []byte(cleartext)
	var encrypted []byte
	switch mm.Construct {
	case mode.AES256_GCM:
		encrypted, err = encrypt(data, key)
	case mode.NACL_SECRETBOX:
		encrypted, err = encryptSecretbox(data, key)
	case mode.XCHACHA20_POLY1305:
		encrypted, err = encryptXChaPo(data, key)
	default:
		return "", errors.New("selected construct not implemented")
	}

	var encodedSalt string
	var encodedCiphertext string
	switch mm.Encoding {
	case mode.BASE32:
		encodedSalt = encode(salt)
		encodedCiphertext = encode(encrypted)
	case mode.BASE62:
		encodedSalt = encodeBase62(salt)
		encodedCiphertext = encodeBase62(encrypted)
	case mode.BASE64:
		encodedSalt = encodeBase64(salt)
		encodedCiphertext = encodeBase64(encrypted)
	case mode.BASE64URL:
		encodedSalt = encodeBase64U(salt)
		encodedCiphertext = encodeBase64U(encrypted)
	default:
		return "", errors.New("selected encoding not implemented")
	}

	m := Message{
		Prefix:     prefix,
		Mode:       *mm,
		Salt:       encodedSalt,
		Ciphertext: encodedCiphertext,
	}
	messageText = m.Text()
	return messageText, err
}

func Decrypt(messageText string, password string) (cleartext string, modeText string, err error) {
	var message *Message
	message, err = NewMessage(messageText)
	if err != nil {
		return "", "", err
	}
	mm := message.Mode

	var decodedSalt []byte
	var decodedCiphertext []byte
	switch mm.Encoding {
	case mode.BASE32:
		decodedSalt, err = decode(message.Salt)
		decodedCiphertext, err = decode(message.Ciphertext)
	case mode.BASE62:
		decodedSalt, err = decodeBase62(message.Salt)
		decodedCiphertext, err = decodeBase62(message.Ciphertext)
	case mode.BASE64:
		decodedSalt, err = decodeBase64(message.Salt)
		decodedCiphertext, err = decodeBase64(message.Ciphertext)
	case mode.BASE64URL:
		decodedSalt, err = decodeBase64U(message.Salt)
		decodedCiphertext, err = decodeBase64U(message.Ciphertext)
	}
	if err != nil {
		return "", "", err
	}

	pw := []byte(password)
	var key *key
	switch mm.Kdf {
	case mode.PBKDF2_SHA256_I10K:
		key = deriveKey(pw, decodedSalt)
	case mode.ARGON2ID_T1_M65536_C4:
		key = deriveKeyArgon2id(pw, decodedSalt)
	case mode.SCRYPT_N32768_R8_P1:
		key = deriveKeyScrypt(pw, decodedSalt)
	default:
		return "", "", errors.New("selected key derivation function not implemented")
	}

	var decrypted []byte
	switch mm.Construct {
	case mode.AES256_GCM:
		decrypted, err = decrypt(decodedCiphertext, key)
	case mode.NACL_SECRETBOX:
		decrypted, err = decryptSecretbox(decodedCiphertext, key)
	case mode.XCHACHA20_POLY1305:
		decrypted, err = decryptXChaPo(decodedCiphertext, key)
	default:
		return "", "", errors.New("selected construct not implemented")
	}

	return string(decrypted), mm.Text(), err
}
