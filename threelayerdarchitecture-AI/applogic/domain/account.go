package domain

import "errors"

type Currency string

const (
	CurrencyUSD Currency = "USD"
)

var (
	ErrAccountIDRequired    = errors.New("account: id is required")
	ErrAccountUserRequired  = errors.New("account: user id is required")
	ErrAccountCurrencyEmpty = errors.New("account: currency is required")
	ErrAmountNegative       = errors.New("account: amount must be greater than zero")
	ErrInsufficientBalance  = errors.New("account: insufficient balance")
)

type Account struct {
	ID       string
	UserID   string
	Balance  int64
	Currency Currency
}

func (a Account) Validate() error {
	switch {
	case a.ID == "":
		return ErrAccountIDRequired
	case a.UserID == "":
		return ErrAccountUserRequired
	case a.Currency == "":
		return ErrAccountCurrencyEmpty
	default:
		return nil
	}
}

func (a Account) Debit(amount int64) (Account, error) {
	if amount <= 0 {
		return Account{}, ErrAmountNegative
	}
	if amount > a.Balance {
		return Account{}, ErrInsufficientBalance
	}
	updated := a
	updated.Balance -= amount
	return updated, nil
}

func (a Account) Credit(amount int64) (Account, error) {
	if amount <= 0 {
		return Account{}, ErrAmountNegative
	}
	updated := a
	updated.Balance += amount
	return updated, nil
}
