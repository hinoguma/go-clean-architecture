package repository

import (
	"context"
	"fmt"

	"app/threelayerdarchitecture/applogic/domain"
	"app/threelayerdarchitecture/appinfra/memorydb"
)

type AccountPort interface {
	Find(ctx context.Context, id string) (domain.Account, error)
	Store(ctx context.Context, account domain.Account) error
}

type TransactionPort interface {
	Record(ctx context.Context, txn domain.Transaction) (string, error)
}

type accountRepository struct {
	client *memorydb.AccountClient
}

func NewAccountRepository(client *memorydb.AccountClient) AccountPort {
	return &accountRepository{client: client}
}

func (r *accountRepository) Find(ctx context.Context, id string) (domain.Account, error) {
	_ = ctx
	record, err := r.client.Get(id)
	if err != nil {
		return domain.Account{}, fmt.Errorf("account repository find: %w", err)
	}
	account := domain.Account{
		ID:       record.ID,
		UserID:   record.UserID,
		Balance:  record.Balance,
		Currency: domain.Currency(record.Currency),
	}
	return account, nil
}

func (r *accountRepository) Store(ctx context.Context, account domain.Account) error {
	_ = ctx
	if err := account.Validate(); err != nil {
		return err
	}
	record := memorydb.AccountRecord{
		ID:       account.ID,
		UserID:   account.UserID,
		Balance:  account.Balance,
		Currency: string(account.Currency),
	}
	if err := r.client.Put(record); err != nil {
		return fmt.Errorf("account repository store: %w", err)
	}
	return nil
}

type transactionRecorder struct {
	client *memorydb.TransactionClient
}

func NewTransactionRecorder(client *memorydb.TransactionClient) TransactionPort {
	return &transactionRecorder{client: client}
}

func (r *transactionRecorder) Record(ctx context.Context, txn domain.Transaction) (string, error) {
	_ = ctx
	record := memorydb.TransactionRecord{
		ID:               txn.ID,
		FromAccountID:    txn.FromAccountID,
		ToAccountID:      txn.ToAccountID,
		Amount:           txn.Amount,
		Currency:         string(txn.Currency),
		OccurredAtUnix:   txn.OccurredAt.Unix(),
		CorrelationID:    txn.CorrelationID,
		InitiatingUserID: txn.InitiatingUserID,
	}
	id, err := r.client.Put(record)
	if err != nil {
		return "", fmt.Errorf("transaction recorder: %w", err)
	}
	return id, nil
}
