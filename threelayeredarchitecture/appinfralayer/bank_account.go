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

func (item *BankAccountRawData) SetByRows(rows *sql.Rows) error {
	return rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
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
	return &BankAccountRepository{
		sqlClient: sqlClient,
	}
}

func (repo *BankAccountRepository) Get(ctx context.Context, id string) (BankAccountRawData, error) {

	stmt, params := buildSelectQueryWithId(id, TableBankAccounts)
	rows, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	if err != nil {
		return BankAccountRawData{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var account BankAccountRawData
		if err := account.SetByRows(); err != nil {
			return BankAccountRawData{}, err
		}
		// first row found and return
		return account, nil
	}

	// no rows found
	return BankAccountRawData{}, sql.ErrNoRows
}


func (repo *BankAccountRepository) Create(ctx context.Context, item BankAccountRawData) error {
	stmt, params := buildInsertQuery(TableBankAccounts, map[string]any{
		"id":         item.ID,
		"created_at": item.CreatedAt,
		"updated_at": item.UpdatedAt,
	},)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	return err
}

func (repo *BankAccountRepository) Update(ctx context.Context, id string, updateFields UpdateFieldRequests) error {
	stmt, params := buildUpdateQueryWithId(id, TableBankAccounts, updateFields)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	return err
}

func (repo *BankAccountRepository) Delete(ctx context.Context, id string) error {
	stmt, params := buildDeleteQueryWithId(id, TableBankAccounts)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	return err
}