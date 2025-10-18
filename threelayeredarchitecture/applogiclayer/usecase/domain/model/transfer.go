package model

type TransactionRecord struct {
	ID                string
	FromBankAccountID string
	ToBankAccountID   string
	Money             Money
	DataItem
}

type TransferRequest struct {
	FromBankAccountID string
	ToBankAccountID   string
	Money             Money
}

type TransferResult struct {
	TransactionRecord TransactionRecord
}
