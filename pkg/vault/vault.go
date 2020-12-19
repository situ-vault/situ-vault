package vault

func Encrypt(clearText string, password string) (string, error) {
	salt := newSalt()
	key := deriveKey([]byte(password), salt)
	encrypted, err := encrypt([]byte(clearText), key)
	cipherText := buildEnvelope(encode(salt), encode(encrypted))
	return cipherText, err
}

func Decrypt(cipherText string, password string) (string, error) {
	salt, encrypted := openEnvelope(cipherText)
	key := deriveKey([]byte(password), decode(salt))
	decrypted, err := decrypt(decode(encrypted), key)
	return string(decrypted), err
}
