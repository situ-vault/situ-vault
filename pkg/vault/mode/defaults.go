package mode

type defaults struct {
	Conservative Mode
	Modern       Mode
	Secretbox    Mode
	XChaCha      Mode
}

var d = defaults{
	Conservative: Mode{
		Construct: AES256_GCM,
		Kdf:       PBKDF2_SHA256_I10K,
		Salt:      R8B,
		Encoding:  BASE32,
	},
	Modern: Mode{
		Construct: AES256_GCM,
		Kdf:       ARGON2ID_T1_M65536_C4,
		Salt:      R16B,
		Encoding:  BASE62,
	},
	Secretbox: Mode{
		Construct: NACL_SECRETBOX,
		Kdf:       SCRYPT_N32768_R8_P1,
		Salt:      R16B,
		Encoding:  BASE64,
	},
	XChaCha: Mode{
		Construct: XCHACHA20_POLY1305,
		Kdf:       SCRYPT_N32768_R8_P1,
		Salt:      R16B,
		Encoding:  BASE64URL,
	},
}

func Defaults() defaults {
	return d
}
