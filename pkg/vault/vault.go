package vault

func Encrypt(clearText string, password string) (string, error) {
	salt := newSalt()
	key := deriveKey([]byte(password), salt)
	encrypted, err := encrypt([]byte(clearText), key)
	cipherText := buildEnvelope(encode(salt), encode(encrypted))
	return cipherText, err
}

func Decrypt(cipherText string, password string) (string, error) {
	salt, encrypted, err := openEnvelope(cipherText)
	if err != nil {
		return "", err
	}
	s, err := decode(salt)
	if err != nil {
		return "", err
	}
	key := deriveKey([]byte(password), s)
	e, err := decode(encrypted)
	if err != nil {
		return "", err
	}
	decrypted, err := decrypt(e, key)
	return string(decrypted), err
}
