package appinfralayer

import (
	"app/experimentarchitecture/applogiclayer/domain/model"
	"app/experimentarchitecture/crosscutting/logger"
	utils2 "app/experimentarchitecture/crosscutting/utils"
	"context"
	"database/sql"
)

var commonLogger = logger.NewStdLogger()

func beforeQueryLog(ctx context.Context, funcName string, query string, args ...any) {
	// get function name to call this
	commonLogger.Info(ctx, utils2.NewLogRequest("query log").
		AdditionalInfo("function", funcName).
		AdditionalInfo("query", query).
		AdditionalInfo("args", args).ToValue(),
	)
}

type sqlConfig struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}

type SqlClient struct {
	db *sql.DB
}

func NewSqlClient(cfg sqlConfig) (SqlClient, error) {
	dsn := cfg.user + ":" + cfg.password + "@tcp(" + cfg.host + ":" + utils2.IntToStr(cfg.port) + ")/" + cfg.dbName + "?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return SqlClient{}, utils2.ErrWrap(err, utils2.NewError("failed to open sql db"))
	}
	return SqlClient{db: db}, nil
}

func (client SqlClient) BeginTx(ctx context.Context, opt model.TxOptions) (*sql.Tx, error) {
	beforeQueryLog(ctx, "BeginTx", "BEGIN")
	rawOpt := sql.TxOptions{}
	tx, err := client.db.BeginTx(ctx, &rawOpt)
	if err != nil {
		return nil, utils2.ErrWrap(err, utils2.NewError("failed to begin tx"))
	}
	return tx, nil
}

type SqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRow(query string, args ...any) *sql.Row
}

func ExecContext(exec SqlExecutor, ctx context.Context, query string, args ...any) (sql.Result, error) {
	beforeQueryLog(ctx, "ExecContext", query, args)
	return exec.ExecContext(ctx, query, args...)
}
func Exec(exec SqlExecutor, query string, args ...any) (sql.Result, error) {
	beforeQueryLog(nil, "Exec", query, args)
	return exec.Exec(query, args...)
}
func QueryContext(exec SqlExecutor, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	beforeQueryLog(ctx, "QueryContext", query, args)
	return exec.QueryContext(ctx, query, args...)
}
func Query(exec SqlExecutor, query string, args ...any) (*sql.Rows, error) {
	beforeQueryLog(nil, "Query", query, args)
	return exec.Query(query, args...)
}
func QueryRowContext(exec SqlExecutor, ctx context.Context, query string, args ...any) *sql.Row {
	beforeQueryLog(ctx, "QueryRowContext", query, args)
	return exec.QueryRowContext(ctx, query, args...)
}
func QueryRow(exec SqlExecutor, query string, args ...any) *sql.Row {
	beforeQueryLog(nil, "QueryRow", query, args)
	return exec.QueryRow(query, args...)
}
