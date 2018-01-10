package utils

import (
	"golang.org/x/crypto/blake2s"

	"golang.org/x/crypto/sha3"
)

//SimpleHashWrap - wrapper for 3rd party library functions
type SimpleHashWrap struct {
}

//Blake2s wrapper method
func (t SimpleHashWrap) Blake2s(data []byte) [32]byte {
	return blake2s.Sum256(data)
}

//Sha256 wrapper method
func (t SimpleHashWrap) Sha256(data []byte) [32]byte {
	return sha3.Sum256(data)
}

//ISimpleHash Interface for SimpleHash module
type ISimpleHash interface {
	Sum256(data []byte) []byte
}

//SimpleHash module
type SimpleHash struct {
	wrapper iSimpleHashWrap
}

//NewSimpleHash - SimpleHash module constructor
func NewSimpleHash(wrapper iSimpleHashWrap) *SimpleHash {
	x12 := new(SimpleHash)
	x12.wrapper = wrapper

	return x12
}

//Sum256 - calculate SHA3 + BLAKE2s hash over byte array
func (t SimpleHash) Sum256(data []byte) []byte {
	sha3Hash := t.wrapper.Sha256(data)
	blakeHash := t.wrapper.Blake2s(sha3Hash[:])
	return blakeHash[:]
}
