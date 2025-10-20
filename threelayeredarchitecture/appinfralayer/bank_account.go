package appinfralayer

import (
	"context"
	"database/sql"
)

const TableBankAccounts = "bank_accounts"

type BankAccountRawData struct {
	ID string
	DatabaseItem
}

func newBankAccountRawDataBySQLRows(rows *sql.Rows) (BankAccountRawData, error) {
	newItem := BankAccountRawData{}
	err := rows.Scan(&newItem.ID, &newItem.CreatedAt, &newItem.UpdatedAt)
	return newItem, err
}

func (item BankAccountRawData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        item.ID,
		"createdAt": item.CreatedAt,
		"updatedAt": item.UpdatedAt,
	}
}

type BankAccountRepositoryIF interface {
	Get(ctx context.Context, id string) (BankAccountRawData, error)
	Create(ctx context.Context, item BankAccountRawData) error
	Update(ctx context.Context, id string, updateFields UpdateFieldRequests) error
	Delete(ctx context.Context, id string) error
}

type BankAccountRepository struct {
	sqlClient SQLClient
}

func NewBankAccountRepository(sqlClient SQLClient) BankAccountRepositoryIF {
	return &TableRepository[BankAccountRawData, string]{
		tablename:   TableBankAccounts,
		sqlClient:   sqlClient,
		convertFunc: newBankAccountRawDataBySQLRows,
	}
}
