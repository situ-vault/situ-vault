package vaultmode

type Salt string

const (
	R8b  Salt = "R8B"  // Random 8 bytes
	R16b Salt = "R16B" // Random 16 bytes
)

type salts struct {
	R8b  Salt
	R16b Salt
}

// intentionally returns private struct
func Salts() salts {
	return salts{
		R8b:  R8b,
		R16b: R16b,
	}
}
