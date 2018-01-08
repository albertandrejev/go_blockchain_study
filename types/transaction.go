package types

//Transaction - contains information about transaction
type Transaction struct {
	TxID string
	Sign string
	Data TransactionData
}
