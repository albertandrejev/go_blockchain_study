package merkle

//Tree - merkle tree structure
type Tree struct {
}

//Init - initialize merkle tree
func (t Tree) Init() []byte {

}

//CalcMerkleNode - calculate hash of two neighbour nodes
func CalcMerkleNode(nodeA []byte, nodeB []byte) []byte {
	nodeSum := append(nodeA, nodeB...)
	return SimpleHash(nodeSum)
}
