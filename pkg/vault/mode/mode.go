package mode

// old:
// "SITU_VAULT_V1##AES256_GCM_PBKDF2_SHA256_ITER10K_SALT8_BASE32##TNSIVLVV6EOGI===##GRDENILPW24R4YDA2I6MKT6JPLG5GM2HWC5S2PR7"

// new proposal:
// "SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##TNSIVLVV6EOGI===##GRDENILPW24R4YDA2I6MKT6JPLG5GM2HWC5S2PR7"
// WIP: "SITU_VAULT_V1##C:AES256_GCM#KDF:ARGON2ID_T1_M65536_C4#SALT:R8B#ENC:BASE64URL##YW55IGNhc===##jA0EAw_BBPOQhPfTDInn-94hXmnBr9D8-4x5"
// WIP: "SITU_VAULT_V1##C:NACL_SECRETBOX#KDF:SYCRYPT_N32768_R8_P1#SALT:R16B#ENC:BASE64##YW55IGNhcm5hbCBwbGV===##jA0EAw/BBPOQhPfTDInn+94hXmnBr9D8+4x5"

type Mode struct {
	Construct Construct             `code:"C"`
	Kdf       KeyDerivationFunction `code:"KDF"`
	Salt      Salt                  `code:"SALT"`
	Encoding  Encoding              `code:"ENC"`
}

type Construct string

const (
	AES256_GCM     Construct = "AES256_GCM"     // AEAD; Standard: Nonce: 12 byte, Tag: 16 byte
	NACL_SECRETBOX Construct = "NACL_SECRETBOX" // XSalsa20 and Poly1305 MAC; Standard: Nonce: 24 byte, Tag: 16 byte
)

type KeyDerivationFunction string

const (
	PBKDF2_SHA256_I10K    KeyDerivationFunction = "PBKDF2_SHA256_I10K"    // 10000 iterations, OpenSSL default
	ARGON2ID_T1_M65536_C4 KeyDerivationFunction = "ARGON2ID_T1_M65536_C4" // parameters as of RFC
	SYCRYPT_N32768_R8_P1  KeyDerivationFunction = "SYCRYPT_N32768_R8_P1"  // parameters as RFC with bigger N
)

type Salt string

const (
	R8B  Salt = "R8B"  // Random 8 bytes
	R16B Salt = "R16B" // Random 16 bytes
)

type Encoding string

const (
	NONE      Encoding = "NONE"      // No encoding, just bytes
	BASE32    Encoding = "BASE32"    // Base32
	BASE64    Encoding = "BASE64"    // Base64
	BASE64URL Encoding = "BASE64URL" // Base64 (URL safe variant)
)
