package utils

type iSimpleHashWrap interface {
	Blake2s(data []byte) [32]byte
	Sha256(data []byte) [32]byte
}
