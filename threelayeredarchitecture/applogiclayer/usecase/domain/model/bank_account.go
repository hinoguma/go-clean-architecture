package model


type BankAccount struct {
	ID      string
	OwnerBankCustomerID string
	Balance Money

	DataItem
}
