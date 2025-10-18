package domain

import (
	"errors"
	"time"
)

var (
	ErrTransactionSameAccount = errors.New("transaction: from and to accounts must differ")
	ErrTransactionAmountZero  = errors.New("transaction: amount must be positive")
)

type Transaction struct {
	ID                string
	FromAccountID     string
	ToAccountID       string
	Amount            int64
	Currency          Currency
	OccurredAt        time.Time
	CorrelationID     string
	InitiatingUserID  string
	CompletionComment string
}

func NewTransaction(from, to Account, amount int64, occurredAt time.Time) (Transaction, error) {
	if err := from.Validate(); err != nil {
		return Transaction{}, err
	}
	if err := to.Validate(); err != nil {
		return Transaction{}, err
	}
	if from.ID == to.ID {
		return Transaction{}, ErrTransactionSameAccount
	}
	if amount <= 0 {
		return Transaction{}, ErrTransactionAmountZero
	}
	return Transaction{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        amount,
		Currency:      from.Currency,
		OccurredAt:    occurredAt,
	}, nil
}
