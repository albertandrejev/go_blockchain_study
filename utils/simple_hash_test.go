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
	crypt := new(simpleHashWrap)

	hash := crypt.Blake2s(blakeData)

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

	hash := simpleHashIntern(simpleInputData, hashMock)

	if bytes.Compare(simpleHashData[:], hash[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, hashData)
	}
}

func TestSimpleHash(t *testing.T) {
	expected := []byte{196, 222, 198, 36, 211, 67, 243, 236, 115, 195, 31, 6, 82, 220,
		236, 49, 171, 43, 31, 145, 153, 220, 161, 107, 210, 207, 33, 156, 134, 169, 94, 173}

	hash := SimpleHash(simpleInputData)

	if bytes.Compare(hash[:], expected[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, expected)
	}
}
