package transaction

import "errors"

var (
	ErrorCancelledContext = errors.New("cancelled context was provided")
	ErrorCommitTx = errors.New("error occurs while committing transaction")
)