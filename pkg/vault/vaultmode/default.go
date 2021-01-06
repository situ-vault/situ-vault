package vaultmode

type defaults struct {
	Conservative Mode
	Modern       Mode
	Secretbox    Mode
	XChaCha      Mode
}

// intentionally returns private struct
func Defaults() defaults {
	return defaults{
		Conservative: Mode{
			Construct: Constructs.Aes256gcm,
			Kdf:       KeyDerivationFunctions.Pbkdf2_sha256_i10k,
			Salt:      Salts.R8b,
			Encoding:  Encodings.Base32,
			Linebreak: Linebreaks.No,
		},
		Modern: Mode{
			Construct: Constructs.Aes256gcm,
			Kdf:       KeyDerivationFunctions.Argon2id_t1_m65536_c4,
			Salt:      Salts.R16b,
			Encoding:  Encodings.Base62,
			Linebreak: Linebreaks.Ch120,
		},
		Secretbox: Mode{
			Construct: Constructs.NaclSecretbox,
			Kdf:       KeyDerivationFunctions.Scrypt_n32768_r8_p1,
			Salt:      Salts.R16b,
			Encoding:  Encodings.Base64,
			Linebreak: Linebreaks.No,
		},
		XChaCha: Mode{
			Construct: Constructs.XChaCha20poly1305,
			Kdf:       KeyDerivationFunctions.Scrypt_n32768_r8_p1,
			Salt:      Salts.R16b,
			Encoding:  Encodings.Base64url,
			Linebreak: Linebreaks.No,
		},
	}
}
