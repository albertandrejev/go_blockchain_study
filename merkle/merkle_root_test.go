package merkle

import (
	"bytes"
	"testing"

	factory_mock "bitbucket.org/albert_andrejev/free_info/factory/mock"
	utils_mock "bitbucket.org/albert_andrejev/free_info/utils/mock"
	"github.com/golang/mock/gomock"
)

var simpleDataToHash = []byte{'1', '2', '3', '4'}
var nodeAHash = []byte{0, 1, 2, 3, 4, 5, 6}
var nodeBHash = []byte{7, 8, 9, 10, 11, 12, 13}

func factoryMockInit(ctrl *gomock.Controller) *factory_mock.MockIMainFactory {
	hashMock := utils_mock.NewMockISimpleHash(ctrl)
	nodeSum := append(nodeAHash, nodeBHash...)
	hashMock.EXPECT().Sum256(gomock.Eq(nodeSum)).Return(simpleDataToHash)

	factoryMock := factory_mock.NewMockIMainFactory(ctrl)
	factoryMock.EXPECT().GetSimpleHash().Return(hashMock)

	return factoryMock
}

func TestSimpleHash_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factoryMock := factoryMockInit(ctrl)

	merkle := NewTree(factoryMock)

	hash := merkle.CalcNode(nodeAHash, nodeBHash)

	if bytes.Compare(simpleDataToHash[:], hash[:]) != 0 {
		t.Errorf("Wrong return hash value.\nActual: %v.\nExpected: %v", hash, simpleDataToHash)
	}
}
