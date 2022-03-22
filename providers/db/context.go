package db

import (
	"context"
	"database/sql"
)

type txKey struct{}

// TxFromContext returns context bounded db transaction
func TxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}

// NewTxContext creates context with db transaction
func NewTxContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
