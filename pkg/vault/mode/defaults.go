package mode

var DefaultConservative = &Mode{
	Construct: AES256_GCM,
	Kdf:       PBKDF2_SHA256_I10K,
	Salt:      R8B,
	Encoding:  BASE32,
}
