package testdata

type Vector struct {
	Password   string
	Ciphertext string
	Cleartext  string
}

func PredefinedDecrypt() Vector {
	return Vector{
		Password:   "test-pw",
		Ciphertext: "SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##TNSIVLVV6EOGI===##GRDENILPW24R4YDA2I6MKT6JPLG5GM2HWC5S2PR7##END",
		Cleartext:  "test-data",
	}
}
