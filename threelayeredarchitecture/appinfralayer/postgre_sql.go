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

func (pg *PostgreSQLClient) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return SQLQueryContext(ctx, pg.client.QueryContext, query, args...)
}


