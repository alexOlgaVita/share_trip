package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func tx[T interface{}](
	ctx context.Context,
	pool *pgxpool.Pool,
	block func(tx pgx.Tx) (*T, error),
) (*T, error) {
	txBegin, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(txBegin pgx.Tx, ctx context.Context) {
		err := txBegin.Rollback(ctx)
		if err != nil {
			return
		}
	}(txBegin, ctx)

	res, err := block(txBegin)
	if err != nil {
		return nil, fmt.Errorf("repoTrip.transacction block: %w", err)
	}

	if err = txBegin.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return res, nil
}
