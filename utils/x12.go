package utils

import (
	"github.com/blockcypher/gox11hash"
	"golang.org/x/crypto/scrypt"
)

//X12HashWrapper - wrapper for 3rd party library functions
type X12HashWrapper struct {
}

//X11 wrapper method
func (t X12HashWrapper) X11(data []byte) []byte {
	return gox11hash.Sum(data)
}

//Scrypt wrapper method
func (t X12HashWrapper) Scrypt(data []byte, salt []byte, N int, r int, p int, keyLen int) ([]byte, error) {
	return scrypt.Key(data, salt, N, r, p, keyLen)
}

//IX12Hash Interface for SimpleHash module
type IX12Hash interface {
	Sum256(data []byte) ([]byte, error)
}

//X12Hash module
type X12Hash struct {
	wrapper iX12HashWrap
}

//NewX12Hash - X12Hash module constructor
func NewX12Hash(wrapper iX12HashWrap) *X12Hash {
	x12 := new(X12Hash)
	x12.wrapper = wrapper

	return x12
}

//Sum256 - x11 hash plus scrypt hash
func (t X12Hash) Sum256(data []byte) ([]byte, error) {
	val := t.wrapper.X11(data)
	scryptHash, err := t.wrapper.Scrypt(val, nil, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return scryptHash, nil
}
