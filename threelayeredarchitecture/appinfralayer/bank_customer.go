package appinfralayer

import (
	"context"
	"database/sql"
)

type BankCustomerRepositoryIF interface {
	Get(ctx context.Context, id string) (BankCustomerRawData, error)
	Create(ctx context.Context, item BankCustomerRawData) error
	Update(ctx context.Context, id string, updateFields UpdateFieldRequests) error
	Delete(ctx context.Context, id string) error
}

const TableBankCustomers = "bank_customers"

type BankCustomerRawData struct {
	ID string
	DatabaseItem
}

func (item BankCustomerRawData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        item.ID,
		"createdAt": item.CreatedAt,
		"updatedAt": item.UpdatedAt,
	}
}

func newBankCustomerRawDataBySQLRows(rows *sql.Rows) (BankCustomerRawData, error) {
	item := BankCustomerRawData{}
	err := rows.Scan(item.ID, item.CreatedAt, item.UpdatedAt)
	return item, err
}

type BankCustomerRepository struct {
	sqlClient SQLClient
}

func NewBankCustomerRepository(sqlClient SQLClient) BankCustomerRepositoryIF {
	return &TableRepository[BankCustomerRawData, string]{
		tablename:   TableBankCustomers,
		sqlClient:   sqlClient,
		convertFunc: newBankCustomerRawDataBySQLRows,
	}
}
