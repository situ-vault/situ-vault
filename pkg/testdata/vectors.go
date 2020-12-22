package testdata

type Vector struct {
	Password   string
	Ciphertext string
	Cleartext  string
}

func PredefinedDecrypt() Vector {
	return Vector{
		Password:   "test-pw",
		Ciphertext: "SITU_VAULT_V1##AES256_GCM_PBKDF2_SHA256_ITER10K_SALT8_BASE32##TNSIVLVV6EOGI===##GRDENILPW24R4YDA2I6MKT6JPLG5GM2HWC5S2PR7",
		Cleartext:  "test-data",
	}
}
