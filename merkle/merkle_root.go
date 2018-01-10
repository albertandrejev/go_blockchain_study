package merkle

import (
	"bitbucket.org/albert_andrejev/free_info/factory"
	"bitbucket.org/albert_andrejev/free_info/types"
)

//Tree - merkle tree structure
type Tree struct {
	factory factory.IMainFactory
}

//NewTree Merkle tree constructor
func NewTree(factory factory.IMainFactory) *Tree {
	tree := new(Tree)
	tree.factory = factory
	return tree
}

//Init - calcualte initial tree leafs from transactions
func (t Tree) Init(transactions []*types.TransactionData) {

}

//CalcNode - calculate hash of two neighbour nodes
func (t Tree) CalcNode(nodeAHash []byte, nodeBHash []byte) []byte {
	sum := t.factory.GetSimpleHash()
	nodeSum := append(nodeAHash, nodeBHash...)
	return sum.Sum256(nodeSum)
}
