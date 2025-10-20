package appinfralayer

import (
	"context"
	"database/sql"
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

type TableRepository[T SQLDatabaseItem[T], idType string | int64] struct {
	tablename   string
	sqlClient   SQLClient
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

func (repo TableRepository[T, idType]) Update(ctx context.Context, id string, updateFields UpdateFieldRequests) error {
	stmt, params := buildUpdateQueryWithId(id, repo.tablename, updateFields)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	return err
}

func (repo TableRepository[T, idType]) Delete(ctx context.Context, id string) error {
	stmt, params := buildDeleteQueryWithId(id, repo.tablename)
	_, err := repo.sqlClient.QueryContext(ctx, stmt, params...)
	return err
}
