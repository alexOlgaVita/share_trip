package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"job4j.ru/share-trip/internal/dto"
	"job4j.ru/share-trip/internal/repository"
)

type TripUsecase struct {
	TripRepo *repository.RepoPg
}

func (u *TripUsecase) CreateTrip(
	ctx context.Context,
	tx pgx.Tx,
	req dto.CreateTripRequest,
) (*dto.Trip, error) {
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
		return nil, fmt.Errorf("repoTrip.Create: %w", err)
	}

	return trip, nil
}
