package types

//Block - contain information about transactions
type Block struct {
	PrevBlockID  string
	BlockID      string
	Nonce        uint64
	ExtraNonce   uint64
	Transactions []*Transaction
}
