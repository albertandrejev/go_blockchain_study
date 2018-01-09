package utils

import (
	"bytes"
	"testing"

	utils_mock "bitbucket.org/albert_andrejev/free_info/utils/mock"
	"github.com/golang/mock/gomock"
)

var blakeData = []byte{'a', 'b', 'c', 'd'}
var shaData = []byte{'a', 'b', 'c', 'd'}

func TestBlake2sLibCall(t *testing.T) {
	expected := []byte{113, 103, 72, 204, 233, 122, 10, 188, 148, 46, 29, 73, 27, 194,
		81, 2, 245, 182, 255, 113, 238, 98, 168, 106, 189, 96, 90, 108, 64, 18, 1, 105}
	crypt := new(simpleHashWrap)

	hash := crypt.Blake2s(blakeData)

	if bytes.Compare(hash[:], expected[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, expected)
	}
}

func TestSha256LibCall(t *testing.T) {
	expected := []byte{111, 111, 18, 148, 113, 89, 13, 44, 145, 128, 76, 129, 43, 87, 80,
		205, 68, 203, 223, 183, 35, 133, 65, 196, 81, 225, 234, 43, 192, 25, 49, 119}
	crypt := new(simpleHashWrap)

	hash := crypt.Sha256(shaData)

	if bytes.Compare(hash[:], expected[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, expected)
	}
}

var simpleInputData = []byte{'a', 'b', 'c', 'd'}
var blake2Data = [32]byte{'d', 'c', 'b', 'a'}
var simpleHashData = [32]byte{'1', '2', '3', '4'}

func simpleHashMockInit(ctrl *gomock.Controller) *utils_mock.MockiSimpleHashWrap {
	crypt := utils_mock.NewMockiSimpleHashWrap(ctrl)
	crypt.EXPECT().Sha256(gomock.Eq(simpleInputData)).Return(blake2Data)
	crypt.EXPECT().Blake2s(gomock.Eq(blake2Data[:])).Return(simpleHashData)

	return crypt
}

func TestSimpleHash_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hashMock := simpleHashMockInit(ctrl)
	simple := NewSimpleHash(hashMock)

	hash := simple.Sum256(simpleInputData)

	if bytes.Compare(simpleHashData[:], hash[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, hashData)
	}
}
