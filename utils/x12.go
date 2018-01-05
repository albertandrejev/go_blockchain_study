package utils

import (
	"hash"

	"golang.org/x/crypto/ripemd160"

	"github.com/blockcypher/gox11hash"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

type hashCrypt struct {
}

type iHashCrypt interface {
	X11([]byte) []byte
	Scrypt(data []byte, salt []byte, N int, r int, p int, keyLen int) ([]byte, error)
}

func (t hashCrypt) X11(data []byte) []byte {
	return gox11hash.Sum(data)
}

func (t hashCrypt) Scrypt(data []byte, salt []byte, N int, r int, p int, keyLen int) ([]byte, error) {
	return scrypt.Key(data, salt, N, r, p, keyLen)
}

//X12Hash - x11 hash plus scrypt hash
func X12Hash(data []byte) ([]byte, error) {
	crypt := new(hashCrypt)
	return x12HashIntern(data, crypt)
}

func x12HashIntern(data []byte, hashCrypt iHashCrypt) ([]byte, error) {
	val := hashCrypt.X11(data)
	scryptHash, err := hashCrypt.Scrypt(val, nil, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return scryptHash, nil
}

//SimpleHash - calculate SHA3 + RIPEMD160 hash over byte array
func SimpleHash(data []byte) []byte {
	return simpleHashIntern(data, sha3.New256(), ripemd160.New())
}

func simpleHashIntern(data []byte, sha3 hash.Hash, ripemd160 hash.Hash) []byte {
	sha3Hash := sha3.Sum(data)
	return ripemd160.Sum(sha3Hash)
}
