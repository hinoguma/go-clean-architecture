package entities

type BankAccount struct {
	ID        string
	UserID    string
	Balance   Money
	CreatedAt int64
	UpdatedAt int64
}
