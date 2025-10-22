package appinfralayer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PostgreSQLClient struct {
	client *sql.DB
}

func NewPostgreSQLClient(db *sql.DB) SQLClient {
	return &PostgreSQLClient{
		client: db,
	}
}

func (pg PostgreSQLClient) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return SQLQueryContext(ctx, pg.client.QueryContext, query, args...)
}

func (pg PostgreSQLClient) TxQueryContext(ctx context.Context, conn TransactionConnection, query string, args ...any) (*sql.Rows, error) {
	return SQLQueryContext(ctx, conn.tx.QueryContext, query, args...)
}

type TableRepositoryIF[T SQLDatabaseItem[T], idType string | int64] interface {
	Get(ctx context.Context, id idType) (T, error)
	Create(ctx context.Context, item T) error
	Update(ctx context.Context, id idType, updateFields UpdateFieldRequests) error
	Delete(ctx context.Context, id idType) error
	Lock(ctx context.Context, id idType, txId string) (T, error)
	TxCreate(ctx context.Context, item T, txId string) error
	TxUpdate(ctx context.Context, id idType, updateFields UpdateFieldRequests, txId string) error
	TxDelete(ctx context.Context, id idType, txId string) error
}

type TableRepository[T SQLDatabaseItem[T], idType string | int64] struct {
	tablename   string
	sqlClient   SQLClient
	txPool      TxConnectionPoolIF
	convertFunc func(rows *sql.Rows) (T, error)
}

func (repo TableRepository[T, idType]) Get(ctx context.Context, id idType) (T, error) {
	var item T
	stmt, params := buildSelectQueryWithId[idType](id, repo.tablename)
	rows, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	if err != nil {
		return item, err
	}
	defer rows.Close()

	item, err = repo.convertFunc(rows)
	if rows.Next() {
		if err != nil {
			return item, err
		}
		// first row found and return
		return item, nil
	}

	// no rows found
	return item, sql.ErrNoRows
}

func (repo TableRepository[T, idType]) Create(ctx context.Context, item T) error {
	stmt, params := buildInsertQuery(repo.tablename, item.ToMap())
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	return err
}

func (repo TableRepository[T, idType]) Update(ctx context.Context, id idType, updateFields UpdateFieldRequests) error {
	stmt, params := buildUpdateQueryWithId(id, repo.tablename, updateFields)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	return err
}

func (repo TableRepository[T, idType]) Delete(ctx context.Context, id idType) error {
	stmt, params := buildDeleteQueryWithId(id, repo.tablename)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	return err
}

func (repo TableRepository[T, idType]) Lock(ctx context.Context, id idType, txId string) (T, error) {
	var item T
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		return item, errors.New(fmt.Sprintf("Not Found tx connection. id:%s", txId))
	}
	stmt, params := buildSelectQueryWithId[idType](id, repo.tablename)
	rows, err := repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	if err != nil {
		return item, err
	}
	defer rows.Close()

	item, err = repo.convertFunc(rows)
	if rows.Next() {
		if err != nil {
			return item, err
		}
		// first row found and return
		return item, nil
	}

	// no rows found
	return item, sql.ErrNoRows
}

func (repo TableRepository[T, idType]) TxCreate(ctx context.Context, item T, txId string) error {
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		return errors.New(fmt.Sprintf("Not Found tx connection. id:%s", txId))
	}
	stmt, params := buildInsertQuery(repo.tablename, item.ToMap())
	_, err := repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	return err
}

func (repo TableRepository[T, idType]) TxUpdate(ctx context.Context, id idType, updateFields UpdateFieldRequests, txId string) error {
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		return errors.New(fmt.Sprintf("Not Found tx connection. id:%s", txId))
	}
	stmt, params := buildUpdateQueryWithId(id, repo.tablename, updateFields)
	_, err := repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	return err
}

func (repo TableRepository[T, idType]) TxDelete(ctx context.Context, id idType, txId string) error {
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		return errors.New(fmt.Sprintf("Not Found tx connection. id:%s", txId))
	}
	stmt, params := buildDeleteQueryWithId(id, repo.tablename)
	_, err := repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	return err
}
