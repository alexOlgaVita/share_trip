package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"job4j.ru/share-trip/internal/domain"
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
	req CreateTripRequest,
) (CreateTripResponse, error) {
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
		slog.String("driverId", req.DriverID),
	)

	logger.Info("create trip started")

	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*CreateTripResponse, error) {
		txLogger := logger.With(
			slog.String("layer", "transaction"),
		)

		txLogger.Info("transaction started")

		domainResp, err := s.TripUsecase.CreateTrip(ctx, tx, toDomainCreateTripRequest(req))
		if err != nil {
			txLogger.Error(
				"create trip usecase failed",
				slog.Any("error", err),
			)
			return nil, fmt.Errorf("usecase.CreateTrip: %w", err)
		}
		serviceResp := fromDomainCreateTripResponse(domainResp)
		txLogger.Info(
			"transaction completed",
			slog.String("trip_id", serviceResp.ID),
		)

		return &serviceResp, nil
	})

	if err != nil {
		logger.Error(
			"create trip failed",
			slog.Any("error", err),
		)
		return CreateTripResponse{}, fmt.Errorf("failed in transaction: %w", err)
	}

	logger.Info(
		"create trip completed",
		slog.String("trip_id", res.ID),
	)

	return *res, nil
}
