package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j.ru/share-trip/internal/observability/logctx"
	"log/slog"
)

func tx[T interface{}](
	ctx context.Context,
	pool *pgxpool.Pool,
	block func(tx pgx.Tx) (*T, error),
) (*T, error) {
	logger := logctx.Logger(ctx).With(
		slog.String("layer", "transaction"),
	)

	logger.Info("begin transaction")

	txBegin, err := pool.Begin(ctx)
	if err != nil {
		logger.Error(
			"failed to begin transaction",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(txBegin pgx.Tx, ctx context.Context) {
		err := txBegin.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger.Error(
				"rollback transaction failed",
				slog.Any("error", err),
			)
			//			return
		}
	}(txBegin, ctx)

	res, err := block(txBegin)
	if err != nil {
		logger.Error(
			"transaction block failed",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("repoTrip.transaction block: %w", err)
	}

	if err = txBegin.Commit(ctx); err != nil {
		logger.Error(
			"failed to commit transaction",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("commit transaction")

	return res, nil
}
