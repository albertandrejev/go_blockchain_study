package utils

type iSimpleHash interface {
	Blake2s(data []byte) [32]byte
}
