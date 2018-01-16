package merkle

import (
	"bitbucket.org/albert_andrejev/free_info/factory"
	"bitbucket.org/albert_andrejev/free_info/types"
	"bitbucket.org/albert_andrejev/free_info/wrappers"
)

//Tree - merkle tree structure
type Tree struct {
	factory factory.IMainFactory
	json    wrappers.IJSONWrapper
}

//NewTree Merkle tree constructor
func NewTree(factory factory.IMainFactory, json wrappers.IJSONWrapper) *Tree {
	tree := new(Tree)
	tree.factory = factory
	tree.json = json
	return tree
}

//Init - calcualte initial tree leafs from transactions
func (t Tree) Init(transactions []*types.Transaction) [][]byte {
	sum := t.factory.GetSimpleHash()
	result := make([][]byte, len(transactions))
	for i := 0; i < len(transactions); i++ {
		transJSON, err := t.json.Encode(transactions[i])
		if err != nil {
			panic(err)
		}
		result[i] = sum.Sum256(transJSON)
	}

	return result
}

//CalcRoot calc Merkle tree root
func (t Tree) CalcRoot(sums [][]byte) []byte {
	sumsCount := len(sums)
	if sumsCount == 1 {
		return sums[0]
	}
	isEven := sumsCount%2 == 0
	if isEven == false {
		sums = append(sums, sums[sumsCount-1])
		sumsCount++
	}

	result := make([][]byte, sumsCount/2)
	for i := 0; i < sumsCount; i += 2 {
		result[i/2] = t.CalcNode(sums[i], sums[i+1])
	}

	return t.CalcRoot(result)
}

//CalcNode - calculate hash of two neighbour nodes
func (t Tree) CalcNode(nodeAHash []byte, nodeBHash []byte) []byte {
	sum := t.factory.GetSimpleHash()
	nodeSum := append(nodeAHash, nodeBHash...)
	return sum.Sum256(nodeSum)
}
