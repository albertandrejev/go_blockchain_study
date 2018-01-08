package utils

import (
	"bytes"
	"testing"

	utils_mock "bitbucket.org/albert_andrejev/free_info/utils/mock"
	"github.com/golang/mock/gomock"
)

var blakeData = []byte{'a', 'b', 'c', 'd'}

func TestBlake2sLibCall(t *testing.T) {
	expected := []byte{113, 103, 72, 204, 233, 122, 10, 188, 148, 46, 29, 73, 27, 194,
		81, 2, 245, 182, 255, 113, 238, 98, 168, 106, 189, 96, 90, 108, 64, 18, 1, 105}
	crypt := new(hashCrypt)

	hash := crypt.Blake2s(blakeData)

	if bytes.Compare(hash[:], expected[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, expected)
	}
}

var simpleInputData = []byte{'a', 'b', 'c', 'd'}
var blake2Data = []byte{'d', 'c', 'b', 'a'}
var simpleHashData = [32]byte{'1', '2', '3', '4'}

func simpleHashSHA3MockInit(ctrl *gomock.Controller) *utils_mock.MockHash {
	crypt := utils_mock.NewMockHash(ctrl)
	appendedRet := append(simpleInputData, blake2Data...)
	crypt.EXPECT().Sum(gomock.Eq(simpleInputData)).Return(appendedRet)

	return crypt
}

func simpleHashBLAKE2MockInit(ctrl *gomock.Controller) *utils_mock.MockiSimpleHash {
	crypt := utils_mock.NewMockiSimpleHash(ctrl)
	crypt.EXPECT().Blake2s(gomock.Eq(blake2Data)).Return(simpleHashData)

	return crypt
}

func TestSimpleHash_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sha3Mock := simpleHashSHA3MockInit(ctrl)
	blake2Mock := simpleHashBLAKE2MockInit(ctrl)

	hash := simpleHashIntern(simpleInputData, sha3Mock, blake2Mock)

	if bytes.Compare(simpleHashData[:], hash[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, hashData)
	}
}

func TestSimpleHash(t *testing.T) {
	expected := []byte{136, 166, 141, 187, 255, 185, 24, 112, 187, 66, 228, 35, 73, 193,
		76, 155, 233, 80, 77, 86, 28, 12, 31, 47, 156, 213, 10, 149, 253, 253, 247, 99}

	hash, err := SimpleHash(simpleInputData)

	if bytes.Compare(hash[:], expected[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, expected)
	}

	if err != nil {
		t.Error("Should run without errors")
	}
}
