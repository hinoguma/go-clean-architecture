package appinfrainterfacelayer

import (
	"app/experimentarchitecture/applogiclayer/domain/model"
	"context"
)

type Transaction interface {
	GetTransactionID() string
	Commit() error
	Rollback() error
}

type TransactionFactory interface {
	BeginTx(ctx context.Context, opt model.TxOptions) (context.Context, Transaction, error)
}
