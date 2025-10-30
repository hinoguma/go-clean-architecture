package appinfralayer

import (
	"app/threelayeredarchitecture/crosscuttinglayer"
	"context"
	"database/sql"
	"errors"
)

var TxNotFound = errors.New("transaction not found")

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
	Put(ctx context.Context, item T) error
	Update(ctx context.Context, id idType, updateFields UpdateFieldRequests) error
	Delete(ctx context.Context, id idType) error
	Lock(ctx context.Context, id idType, txId string) (T, error)
	TxCreate(ctx context.Context, item T, txId string) error
	TxPut(ctx context.Context, item T, txId string) error
	TxUpdate(ctx context.Context, id idType, updateFields UpdateFieldRequests, txId string) error
	TxDelete(ctx context.Context, id idType, txId string) error
}

type TableRepository[T SQLDatabaseItem[T], idType string | int64] struct {
	tablename   string
	sqlClient   SQLClient
	txPool      TxConnectionPoolIF
	convertFunc func(rows *sql.Rows) (T, error)
}

func (repo TableRepository[T, idType]) Get(ctx context.Context, id idType) (item T, err error) {
	stmt, params := buildSelectQueryWithId[idType](id, repo.tablename)
	rows, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	if err != nil {
		err = crosscuttinglayer.ErrLift(err, ctx)
		return item, err
	}
	defer func() {
		tmpErr := rows.Close()
		if tmpErr != nil {
			err = errors.Join(err, tmpErr)
		}
	}()

	item, err = repo.convertFunc(rows)
	if err != nil {
		err = crosscuttinglayer.ErrLift(err, ctx)
		return item, err
	}
	// success
	err = nil
	return item, err
}

func (repo TableRepository[T, idType]) Create(ctx context.Context, item T) error {
	stmt, params := buildInsertQuery(repo.tablename, item.ToMap())
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	if err != nil {
		return crosscuttinglayer.ErrLift(err, ctx)
	}
	return nil
}

func (repo TableRepository[T, idType]) Put(ctx context.Context, item T) error {
	stmt, params := buildUpsertQuery(repo.tablename, item.ToMap())
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	if err != nil {
		return crosscuttinglayer.ErrLift(err, ctx)
	}
	return nil
}

func (repo TableRepository[T, idType]) Update(ctx context.Context, id idType, updateFields UpdateFieldRequests) error {
	stmt, params := buildUpdateQueryWithId(id, repo.tablename, updateFields)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	if err != nil {
		return crosscuttinglayer.ErrLift(err, ctx)
	}
	return nil
}

func (repo TableRepository[T, idType]) Delete(ctx context.Context, id idType) error {
	stmt, params := buildDeleteQueryWithId(id, repo.tablename)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	if err != nil {
		return crosscuttinglayer.ErrLift(err, ctx)
	}
	return nil
}

func (repo TableRepository[T, idType]) Lock(ctx context.Context, id idType, txId string) (item T, err error) {
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		tmpErr := crosscuttinglayer.NewError(TxNotFound, ctx)
		tmpErr.Attr("txId", txId)
		err = tmpErr
		return item, err
	}
	stmt, params := buildSelectQueryWithId[idType](id, repo.tablename)
	var rows *sql.Rows
	rows, err = repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	if err != nil {
		err = crosscuttinglayer.ErrLift(err, ctx)
		return item, err
	}
	defer func() {
		tmpErr := rows.Close()
		if tmpErr != nil {
			err = errors.Join(err, tmpErr)
		}
	}()

	item, err = repo.convertFunc(rows)
	if err != nil {
		err = crosscuttinglayer.ErrLift(err, ctx)
		return item, err
	}
	err = nil
	return item, err
}

func (repo TableRepository[T, idType]) TxCreate(ctx context.Context, item T, txId string) error {
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		err := crosscuttinglayer.NewError(TxNotFound, ctx)
		err.Attr("txId", txId)
		return err
	}
	stmt, params := buildInsertQuery(repo.tablename, item.ToMap())
	_, err := repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	if err != nil {
		return crosscuttinglayer.ErrLift(err, ctx)
	}
	return nil
}

func (repo TableRepository[T, idType]) TxPut(ctx context.Context, item T, txId string) error {
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		err := crosscuttinglayer.NewError(TxNotFound, ctx)
		err.Attr("txId", txId)
		return err
	}
	stmt, params := buildUpsertQuery(repo.tablename, item.ToMap())
	_, err := repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	if err != nil {
		return crosscuttinglayer.ErrLift(err, ctx)
	}
	return nil
}

func (repo TableRepository[T, idType]) TxUpdate(ctx context.Context, id idType, updateFields UpdateFieldRequests, txId string) error {
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		err := crosscuttinglayer.NewError(TxNotFound, ctx)
		err.Attr("txId", txId)
		return err
	}
	stmt, params := buildUpdateQueryWithId(id, repo.tablename, updateFields)
	_, err := repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	if err != nil {
		return crosscuttinglayer.ErrLift(err, ctx)
	}
	return nil
}

func (repo TableRepository[T, idType]) TxDelete(ctx context.Context, id idType, txId string) error {
	conn, ok := repo.txPool.Get(txId)
	if !ok || !conn.HasTx() {
		err := crosscuttinglayer.NewError(TxNotFound, ctx)
		err.Attr("txId", txId)
		return err
	}
	stmt, params := buildDeleteQueryWithId(id, repo.tablename)
	_, err := repo.sqlClient.TxQueryContext(ctx, conn, stmt, params...)
	if err != nil {
		return crosscuttinglayer.ErrLift(err, ctx)
	}
	return nil
}
