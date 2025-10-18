package domain

import (
	"time"
)

type Clock interface {
	NowUTC() time.Time
}

type defaultClock struct{}

func (defaultClock) NowUTC() time.Time { return time.Now().UTC() }

type TransferService interface {
	Transfer(from, to Account, amount int64) (Account, Account, Transaction, error)
}

type transferService struct {
	clock Clock
}

func NewTransferService(clock Clock) TransferService {
	if clock == nil {
		clock = defaultClock{}
	}
	return transferService{clock: clock}
}

func (s transferService) Transfer(from, to Account, amount int64) (Account, Account, Transaction, error) {
	debited, err := from.Debit(amount)
	if err != nil {
		return Account{}, Account{}, Transaction{}, err
	}
	credited, err := to.Credit(amount)
	if err != nil {
		return Account{}, Account{}, Transaction{}, err
	}
	transaction, err := NewTransaction(debited, credited, amount, s.clock.NowUTC())
	if err != nil {
		return Account{}, Account{}, Transaction{}, err
	}
	return debited, credited, transaction, nil
}
