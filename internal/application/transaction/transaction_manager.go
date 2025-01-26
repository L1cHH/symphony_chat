package transaction

import (
	"context"
	"database/sql"
)

//Interface that implements *sql.DB and *sql.Tx
//So we can use our queries in transaction context or non-transaction context
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

//Interface that implements a transaction manager
type TransactionManager interface {
	WithinTransaction(ctx context.Context, txFunc func(txCtx context.Context) error) error
}

//Type that represents a context key
type ctxKey string

//Context key enum
const (
	TransactionCtxKey ctxKey = "tx"
)

//Function that checks if context has a transaction
func IsTransaction(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value("tx").(*sql.Tx); ok {
		return tx
	}
	return nil
}