package appinfralayer

import (
	"context"
	"database/sql"
)

type TransactionRecordRepositoryIF interface {
	Get(ctx context.Context, id string) (TransactionRecordRawData, error)
	Create(ctx context.Context, item TransactionRecordRawData) error
	Update(ctx context.Context, id string, updateFields UpdateFieldRequests) error
	Delete(ctx context.Context, id string) error

	// With Transaction
	Lock(ctx context.Context, id string, txId string) (TransactionRecordRawData, error)
	TxCreate(ctx context.Context, item TransactionRecordRawData, txId string) error
	TxUpdate(ctx context.Context, id string, updateFields UpdateFieldRequests, txId string) error
	TxDelete(ctx context.Context, id string, txId string) error
}

const TableTransactionRecords = "transaction_records"

type TransactionRecordRawData struct {
	ID string
	DatabaseItem
}

func (item TransactionRecordRawData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        item.ID,
		"createdAt": item.CreatedAt,
		"updatedAt": item.UpdatedAt,
	}
}

func newTransactionRecordRawDataBySQLRows(rows *sql.Rows) (TransactionRecordRawData, error) {
	item := TransactionRecordRawData{}
	err := rows.Scan(item.ID, item.CreatedAt, item.UpdatedAt)
	return item, err
}

type TransactionRecordRepository struct {
	sqlClient SQLClient
}

func NewTransactionRecordRepository(sqlClient SQLClient, txPool TxConnectionPoolIF) TransactionRecordRepositoryIF {
	return &TableRepository[TransactionRecordRawData, string]{
		tablename:   TableTransactionRecords,
		sqlClient:   sqlClient,
		txPool:      txPool,
		convertFunc: newTransactionRecordRawDataBySQLRows,
	}
}
