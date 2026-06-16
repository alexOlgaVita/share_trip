package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"job4j.ru/share-trip/internal/dto"
	"job4j.ru/share-trip/internal/observability/logctx"
	"job4j.ru/share-trip/internal/repository"
	"log/slog"
)

type TripUsecase struct {
	TripRepo *repository.RepoPg
}

func (u *TripUsecase) CreateTrip(
	ctx context.Context,
	tx pgx.Tx,
	req dto.CreateTripRequest,
) (*dto.Trip, error) {
	logger := logctx.Logger(ctx).With(
		slog.String("layer", "usecase"),
		slog.String("usecase", "TripUsecase.CreateTrip"),
		slog.String("client_id", req.DriverId),
	)

	logger.Info("create trip usecase started")

	id := uuid.NewString()
	trip, err := u.TripRepo.Create(ctx, dto.Trip{
		ID:             id,
		DriverId:       req.DriverId,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime,
		AvailableSeats: req.AvailableSeats,
		Status:         dto.TripStatusDraft,
	})
	if err != nil {
		logger.Error(
			"repository create trip failed",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("repoTrip.Create: %w", err)
	}

	logger.Info(
		"create trip usecase completed",
		slog.String("trip_id", trip.ID),
	)

	return trip, nil
}
