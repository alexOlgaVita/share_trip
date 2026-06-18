package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/internal/dto"
	"job4j.ru/share-trip/internal/observability/logctx"
	"job4j.ru/share-trip/internal/observability/metrics"
	"log/slog"
	"time"
)

type TripService struct {
	logger      *slog.Logger
	metrics     *metrics.Metrics
	Pool        *pgxpool.Pool
	TripUsecase *domain.TripUsecase
}

func NewTripService(
	logger *slog.Logger,
	metrics *metrics.Metrics,
	Pool *pgxpool.Pool,
	TripUsecase *domain.TripUsecase,
) *TripService {
	return &TripService{
		logger:      logger,
		metrics:     metrics,
		Pool:        Pool,
		TripUsecase: TripUsecase,
	}
}

func (s *TripService) CreateTrip(
	ctx context.Context,
	req dto.CreateTripRequest,
) (*dto.Trip, error) {
	ctx, span := otel.Tracer("TripService").Start(ctx, "TripService.CreateTrip")
	defer span.End()

	started := time.Now()
	result := "success"

	defer func() {
		s.metrics.TripCreateTotal.WithLabelValues(result).Inc()
		s.metrics.TripCreateDuration.WithLabelValues(result).
			Observe(time.Since(started).Seconds())
	}()

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
