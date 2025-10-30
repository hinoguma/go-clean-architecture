package appinfralayer

import (
	"app/threelayeredarchitecture/crosscuttinglayer"
	"context"
	"database/sql"
	"errors"
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
	TxQueryContext(ctx context.Context, conn TransactionConnection, query string, args ...any) (*sql.Rows, error)
}

type ExecSQLQueryFunc func(ctx context.Context, query string, args ...any) (*sql.Rows, error)

func SQLQueryContext(ctx context.Context, queryFunc ExecSQLQueryFunc, query string, args ...any) (*sql.Rows, error) {
	// log, metrics, tracing, etc.

	row, err := queryFunc(ctx, query, args...)
	// log, metrics, tracing, etc.

	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			err = crosscuttinglayer.NewDataNotFound(err, ctx)
		} else {
			err = crosscuttinglayer.NewError(err, ctx)
		}
		return row, err
	}

	// success
	return row, nil
}

func buildSelectQueryWithId[T string | int64](id T, tableName string) (string, []any) {
	return buildSelectQueryWithSingleCondition("id", id, tableName)
}

func buildSelectQueryWithSingleCondition[T string | int64](field string, val T, tableName string) (string, []any) {
	query := "SELECT * FROM " + tableName + " WHERE $1 = $2"
	values := []any{field, val}
	return query, values
}

func buildSelectQueryWithIdempotencyKey(key string, tableName string) (string, []any) {
	return buildSelectQueryWithSingleCondition("idempotencyKey", key, tableName)
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

// todo: check if postgre has upsert statement
func buildUpsertQuery(tableName string, fields map[string]any) (string, []any) {
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

func (requests *UpdateFieldRequests) Append(field string, value any) *UpdateFieldRequests {
	*requests = append(*requests, UpdateFieldRequest{
		FieldName: field,
		NewValue:  value,
	})
	return requests
}

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
