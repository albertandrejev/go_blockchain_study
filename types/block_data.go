package types

//BlockData - hashable block data
type BlockData struct {
	PrevBlockID string
	MerkleRoot  string
	Nonce       uint64
	ExtraNonce  string
	Target      uint32
	Timestamp   int64
}
