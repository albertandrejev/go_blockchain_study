package types

//Block - contain information about transactions
type Block struct {
	BlockID      string
	Data         BlockData
	Transactions []*Transaction
}
