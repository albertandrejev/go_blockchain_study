package utils

import (
	"hash"

	"golang.org/x/crypto/blake2s"

	"golang.org/x/crypto/sha3"
)

type simpleHash struct {
}

func (t hashCrypt) Blake2s(data []byte) [32]byte {
	return blake2s.Sum256(data)
}

//SimpleHash - calculate SHA3 + BLAKE2s hash over byte array
func SimpleHash(data []byte) ([]byte, error) {
	crypt := new(hashCrypt)
	return simpleHashIntern(data, sha3.New256(), crypt), nil
}

func simpleHashIntern(data []byte, sha3256 hash.Hash, blake2 iSimpleHash) []byte {
	sha3Hash := sha3256.Sum(data)[len(data):]
	blakeHash := blake2.Blake2s(sha3Hash)
	return blakeHash[:]
}
