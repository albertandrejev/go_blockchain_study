package merkle

import "bitbucket.org/albert_andrejev/free_info/factory"

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

//CalcMerkleNode - calculate hash of two neighbour nodes
func (t Tree) CalcMerkleNode(nodeA []byte, nodeB []byte) []byte {
	nodeSum := append(nodeA, nodeB...)
	return t.factory.GetSimpleHash(nodeSum)
}
