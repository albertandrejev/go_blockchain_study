package utils

import (
	"github.com/blockcypher/gox11hash"
	"golang.org/x/crypto/scrypt"
)

//X12Hash - x11 hash plus scrypt hash
func X12Hash(data []byte) ([]byte, error) {
	val := gox11hash.Sum(data)
	scryptHash, err := scrypt.Key(val, nil, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return scryptHash, nil
}
