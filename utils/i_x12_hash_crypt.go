package utils

type iX12HashCrypt interface {
	X11([]byte) []byte
	Scrypt(data []byte, salt []byte, N int, r int, p int, keyLen int) ([]byte, error)
}
