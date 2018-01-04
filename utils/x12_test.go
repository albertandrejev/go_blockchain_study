package utils

import (
	"bytes"
	"errors"
	"testing"

	mock_utils "bitbucket.org/albert_andrejev/free_info/utils/mock"
	"github.com/golang/mock/gomock"
)

var x11Data = []byte{'a', 'b', 'c', 'd'}
var scryptData = []byte{'d', 'c', 'b', 'a'}
var hashData = []byte{'1', '2', '3', '4'}

func mockInit(ctrl *gomock.Controller, err error) *mock_utils.MockiHashCrypt {
	crypt := mock_utils.NewMockiHashCrypt(ctrl)
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

func TestX12Hash_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	crypt := mockInit(ctrl, nil)

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

	crypt := mockInit(ctrl, err)

	hash, err := x12HashIntern(x11Data, crypt)

	if hash != nil {
		t.Error("Should not contain any value")
	}

	if err == nil {
		t.Error("Should raise an error")
	}
}
