package merkle

import "bitbucket.org/albert_andrejev/free_info/types"

//ITree Merkle tree interface
type ITree interface {
	Init(transactions []*types.TransactionData) [][]byte
	CalcRoot(sums [][]byte) []byte
	CalcNode(nodeAHash []byte, nodeBHash []byte) []byte
}
