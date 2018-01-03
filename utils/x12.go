package utils

import (
	"github.com/blockcypher/gox11hash"
	"golang.org/x/crypto/scrypt"
)

type hashCrypt struct {
}

type iHashCrypt interface {
	X11([]byte) []byte
	Scrypt([]byte) ([]byte, error)
}

func (t hashCrypt) X11(data []byte) []byte {
	return gox11hash.Sum(data)
}

func (t hashCrypt) Scrypt(data []byte) ([]byte, error) {
	return scrypt.Key(data, nil, 32768, 8, 1, 32)
}

//X12Hash - x11 hash plus scrypt hash
func X12Hash(data []byte) ([]byte, error) {
	crypt := new(hashCrypt)
	return x12HashIntern(data, crypt)
}

func x12HashIntern(data []byte, hashCrypt iHashCrypt) ([]byte, error) {
	val := hashCrypt.X11(data)
	scryptHash, err := hashCrypt.Scrypt(val)
	if err != nil {
		return nil, err
	}
	return scryptHash, nil
}
