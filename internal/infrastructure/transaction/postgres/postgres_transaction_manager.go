package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"symphony_chat/internal/application/transaction"
)

type PostgresTransactionManager struct {
	db *sql.DB
}

func NewPostgresTransactionManager(db *sql.DB) *PostgresTransactionManager {
	return &PostgresTransactionManager{
		db: db,
	}
}

func (ptm *PostgresTransactionManager) WithinTransaction(ctx context.Context, txFunc func(txCtx context.Context) error) error {

	tx, err := ptm.db.BeginTx(ctx, nil)
	if err != nil {
		return transaction.ErrorCancelledContext
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, transaction.TransactionCtxKey, tx)

	err = txFunc(txCtx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%w: %w", transaction.ErrorCommitTx, err)
	}
	
	return nil
}
