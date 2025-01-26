package transaction

import "context"

type TransactionManager interface {
	WithinTransaction(ctx context.Context, txFunc func(txCtx context.Context) error) error
}