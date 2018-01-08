package utils

import (
	"golang.org/x/crypto/blake2s"

	"golang.org/x/crypto/sha3"
)

type simpleHashWrap struct {
}

func (t simpleHashWrap) Blake2s(data []byte) [32]byte {
	return blake2s.Sum256(data)
}

func (t simpleHashWrap) Sha256(data []byte) [32]byte {
	return sha3.Sum256(data)
}

//SimpleHash - calculate SHA3 + BLAKE2s hash over byte array
func SimpleHash(data []byte) []byte {
	wrapper := new(simpleHashWrap)
	return simpleHashIntern(data, wrapper)
}

func simpleHashIntern(data []byte, wrapper iSimpleHashWrap) []byte {
	sha3Hash := wrapper.Sha256(data)
	blakeHash := wrapper.Blake2s(sha3Hash[:])
	return blakeHash[:]
}
