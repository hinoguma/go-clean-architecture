package appinfralayer

import (
	"app/threelayeredarchitecture/crosscuttinglayer"
	"context"
	"database/sql"
	"errors"
)

type TransactionRecordRepositoryIF interface {
	Get(ctx context.Context, id string) (TransactionRecordRawData, error)
	GetByIdempotencyKey(ctx context.Context, key string) (TransactionRecordRawData, error)
	Create(ctx context.Context, item TransactionRecordRawData) error
	Put(ctx context.Context, item TransactionRecordRawData) error
	Update(ctx context.Context, id string, updateFields UpdateFieldRequests) error
	Delete(ctx context.Context, id string) error

	// With Transaction
	Lock(ctx context.Context, id string, txId string) (TransactionRecordRawData, error)
	TxCreate(ctx context.Context, item TransactionRecordRawData, txId string) error
	TxPut(ctx context.Context, item TransactionRecordRawData, txId string) error
	TxUpdate(ctx context.Context, id string, updateFields UpdateFieldRequests, txId string) error
	TxDelete(ctx context.Context, id string, txId string) error
}

const TableTransactionRecords = "transaction_records"

type TransactionRecordRawData struct {
	ID                string
	Type              string
	IdempotencyKey    string
	FromBankAccountID string
	ToBankAccountID   string
	Amount            int64
	Currency          string
	TransferStatus    string
	DatabaseItem
}

func (item TransactionRecordRawData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                item.ID,
		"type":              item.Type,
		"idempotencyKey":    item.IdempotencyKey,
		"fromBankAccountId": item.FromBankAccountID,
		"toBankAccountId":   item.ToBankAccountID,
		"amount":            item.Amount,
		"currency":          item.Currency,
		"transferStatus":    item.TransferStatus,
		"createdAt":         item.CreatedAt,
		"updatedAt":         item.UpdatedAt,
	}
}

func newTransactionRecordRawDataBySQLRows(rows *sql.Rows) (TransactionRecordRawData, error) {
	item := TransactionRecordRawData{}
	if rows == nil {
		return item, crosscuttinglayer.NewAppUtilError("rows is nil", context.TODO())
	}
	err := rows.Scan(
		item.ID,
		item.Type,
		item.IdempotencyKey,
		item.FromBankAccountID,
		item.ToBankAccountID,
		item.Amount,
		item.Currency,
		item.CreatedAt,
		item.UpdatedAt,
	)
	return item, err
}

type TransactionRecordRepository struct {
	TableRepository[TransactionRecordRawData, string]
}

func (repo TransactionRecordRepository) GetByIdempotencyKey(ctx context.Context, key string) (item TransactionRecordRawData, err error) {
	stmt, params := buildSelectQueryWithIdempotencyKey(key, repo.tablename)
	var rows *sql.Rows
	rows, err = repo.sqlClient.QueryContext(ctx, stmt, params...)
	if err != nil {
		return item, crosscuttinglayer.ErrLift(err, ctx)
	}
	defer func() {
		tmpErr := rows.Close()
		if tmpErr != nil {
			err = errors.Join(err, tmpErr)
		}
	}()

	item, err = repo.convertFunc(rows)
	if err != nil {
		return item, crosscuttinglayer.ErrLift(err, ctx)
	}

	// success
	err = nil
	return item, err
}

func NewTransactionRecordRepository(sqlClient SQLClient, txPool TxConnectionPoolIF) TransactionRecordRepositoryIF {
	repo := TransactionRecordRepository{}
	repo.tablename = TableTransactionRecords
	repo.sqlClient = sqlClient
	repo.txPool = txPool
	repo.convertFunc = newTransactionRecordRawDataBySQLRows
	return &repo
}
