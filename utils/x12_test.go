package utils

import (
	"bytes"
	"errors"
	"testing"

	utils_mock "bitbucket.org/albert_andrejev/free_info/utils/mock"
	"github.com/golang/mock/gomock"
)

var x11Data = []byte{'a', 'b', 'c', 'd'}
var scryptData = []byte{'d', 'c', 'b', 'a'}
var hashData = []byte{'1', '2', '3', '4'}

func x12MockInit(ctrl *gomock.Controller, err error) *utils_mock.MockiX12HashWrap {
	crypt := utils_mock.NewMockiX12HashWrap(ctrl)
	crypt.EXPECT().X11(gomock.Eq(x11Data)).Return(scryptData)
	crypt.EXPECT().Scrypt(
		gomock.Eq(scryptData),
		gomock.Nil(),
		gomock.Eq(32768),
		gomock.Eq(8),
		gomock.Eq(1),
		gomock.Eq(32),
	).Return(hashData, err)

	return crypt
}

func TestX11LibCall(t *testing.T) {
	expected := []byte{154, 30, 187, 231, 240, 48, 57, 248, 134, 114, 110, 192, 76, 29,
		131, 38, 53, 226, 137, 12, 76, 230, 163, 231, 135, 102, 46, 96, 233, 150, 25, 157}
	crypt := new(x12HashWrapper)

	hash := crypt.X11(x11Data)

	if bytes.Compare(hash[:], expected[:]) != 0 {
		t.Error("wrong return hash value")
	}
}

func TestScryptLibCall(t *testing.T) {
	expected := []byte{62, 237, 239, 118, 106, 64, 235, 243, 104, 241, 73, 33, 224, 58, 228,
		34, 247, 190, 94, 114, 139, 199, 203, 228, 238, 146, 81, 73, 164, 182, 168, 109}
	crypt := new(x12HashWrapper)

	hash, err := crypt.Scrypt(scryptData, nil, 32768, 8, 1, 32)

	if bytes.Compare(hash[:], expected[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, expected)
	}

	if err != nil {
		t.Error("Should run without errors")
	}
}

func TestX12HashMain(t *testing.T) {
	expected := []byte{123, 214, 106, 9, 98, 141, 104, 80, 210, 164, 80, 93, 65, 148,
		40, 175, 239, 84, 130, 143, 158, 43, 247, 239, 211, 90, 127, 12, 5, 29, 70, 254}

	hash, err := X12Hash(x11Data)

	if bytes.Compare(expected[:], hash[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, expected)
	}

	if err != nil {
		t.Error("Should run without errors")
	}
}

func TestX12Hash_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	crypt := x12MockInit(ctrl, nil)

	hash, err := x12HashIntern(x11Data, crypt)

	if bytes.Compare(hashData[:], hash[:]) != 0 {
		t.Error("wrong return hash value")
	}

	if err != nil {
		t.Error("Should run without errors")
	}
}

func TestX12Hash_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	err := errors.New("Unittest error")

	crypt := x12MockInit(ctrl, err)

	hash, err := x12HashIntern(x11Data, crypt)

	if hash != nil {
		t.Error("Should not contain any value")
	}

	if err == nil {
		t.Error("Should raise an error")
	}
}
