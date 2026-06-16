package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/internal/dto"
	"job4j.ru/share-trip/internal/observability/logctx"
	"log/slog"
)

type TripService struct {
	Pool        *pgxpool.Pool
	TripUsecase *domain.TripUsecase
}

func (s *TripService) CreateTrip(
	ctx context.Context,
	req dto.CreateTripRequest,
) (*dto.Trip, error) {
	logger := logctx.Logger(ctx).With(
		slog.String("service", "TripService"),
		slog.String("operation", "CreateTrip"),
		slog.String("driverId", req.DriverId),
	)

	logger.Info("create trip started")

	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.Trip, error) {
		txLogger := logger.With(
			slog.String("layer", "transaction"),
		)

		txLogger.Info("transaction started")

		resp, err := s.TripUsecase.CreateTrip(ctx, tx, dto.CreateTripRequest{
			DriverId:       req.DriverId,
			FromPoint:      req.FromPoint,
			ToPoint:        req.ToPoint,
			DepartureTime:  req.DepartureTime,
			AvailableSeats: req.AvailableSeats,
		})
		if err != nil {
			txLogger.Error(
				"create trip usecase failed",
				slog.Any("error", err),
			)
			return nil, fmt.Errorf("usecase.CreateTrip: %w", err)
		}
		txLogger.Info(
			"transaction completed",
			slog.String("trip_id", resp.ID),
		)

		return resp, nil
	})

	if err != nil {
		logger.Error(
			"create trip failed",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	logger.Info(
		"create trip completed",
		slog.String("trip_id", res.ID),
	)

	return res, nil
}
