package utils

import (
	"github.com/blockcypher/gox11hash"
	"golang.org/x/crypto/scrypt"
)

type x12HashWrapper struct {
}

func (t x12HashWrapper) X11(data []byte) []byte {
	return gox11hash.Sum(data)
}

func (t x12HashWrapper) Scrypt(data []byte, salt []byte, N int, r int, p int, keyLen int) ([]byte, error) {
	return scrypt.Key(data, salt, N, r, p, keyLen)
}

//X12Hash - x11 hash plus scrypt hash
func X12Hash(data []byte) ([]byte, error) {
	wrapper := new(x12HashWrapper)
	return x12HashIntern(data, wrapper)
}

func x12HashIntern(data []byte, wrapper iX12HashWrap) ([]byte, error) {
	val := wrapper.X11(data)
	scryptHash, err := wrapper.Scrypt(val, nil, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return scryptHash, nil
}
