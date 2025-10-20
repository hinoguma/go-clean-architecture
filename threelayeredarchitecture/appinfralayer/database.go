package appinfralayer

import (
	"context"
	"database/sql"
)

type SQLDatabaseItem[T any] interface {
	ToMap() map[string]interface{}
}

type DatabaseItem struct {
	CreatedAt int64
	UpdatedAt int64
}

type SQLClient interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type ExecSQLQueryFunc func(ctx context.Context, query string, args ...any) (*sql.Rows, error)

func SQLQueryContext(ctx context.Context, f ExecSQLQueryFunc, query string, args ...any) (*sql.Rows, error) {
	// log, metrics, tracing, etc.

	row, err := f(ctx, query, args...)

	// log, metrics, tracing, etc.

	return row, err
}

func buildSelectQueryWithId[T string | int64](id T, tableName string) (string, []any) {
	query := "SELECT * FROM " + tableName + " WHERE id = $1"
	values := []any{id}
	return query, values
}

func buildUpdateQueryWithId[T string | int64](id T, tableName string, updateFields UpdateFieldRequests) (string, []any) {
	setClause, values := updateFields.ToSQLSetClause()
	query := "UPDATE " + tableName + " SET " + setClause + " WHERE id = $" + string(len(values)+1)
	values = append(values, id)
	return query, values
}

func buildDeleteQueryWithId[T string | int64](id T, tableName string) (string, []any) {
	query := "DELETE FROM " + tableName + " WHERE id = $1"
	values := []any{id}
	return query, values
}

func buildInsertQuery(tableName string, fields map[string]any) (string, []any) {
	columns := ""
	placeholders := ""
	values := make([]any, 0)
	i := 1
	for col, val := range fields {
		if i > 1 {
			columns += ", "
			placeholders += ", "
		}
		columns += col
		placeholders += "$" + string(i)
		values = append(values, val)
		i++
	}
	query := "INSERT INTO " + tableName + " (" + columns + ") VALUES (" + placeholders + ")"
	return query, values
}

type UpdateFieldRequest struct {
	FieldName string
	NewValue  any
}
type UpdateFieldRequests []UpdateFieldRequest

func (requests UpdateFieldRequests) ToSQLSetClause() (string, []any) {
	setClause := ""
	values := make([]any, 0)
	for i, req := range requests {
		if i > 0 {
			setClause += ", "
		}
		setClause += req.FieldName + " = $" + string(i+1)
		values = append(values, req.NewValue)
	}

	// example: "field1 = $1, field2 = $2", [value1, value2]
	return setClause, values
}
