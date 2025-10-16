package entities

type TransactionRecord struct {
	ID                string
	FromBankAccountID string
	ToBankAccountID   string
	Money             Money
}
