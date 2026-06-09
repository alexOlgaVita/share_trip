package domain

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"job4j.ru/share-trip/internal/dto"
)

func (u *TripUsecase) MoveTripDraftToPublish(
	ctx context.Context,
	tx pgx.Tx,
	req dto.UpdateTripRequest,
) (*dto.Trip, error) {

	trip, err := u.TripRepo.GetForUpdateByID(ctx, tx, req.TripID)
	if err != nil {
		return nil, fmt.Errorf("tripRepository.GetForUpdateByID: %w", err)
	}

	if trip.DriverId != req.ClientID {
		return nil, fmt.Errorf("forbidden: client %s is not driver of trip %s", req.ClientID, req.TripID)
	}

	if trip.Status == dto.TripStatusPublished {
		return &trip, nil
	}

	if trip.Status != dto.TripStatusDraft {
		return nil, fmt.Errorf("invalid trip status: expected %s, got %s", dto.TripStatusDraft, trip.Status)
	}

	err = u.TripRepo.UpdateStatus(ctx, tx, trip.ID, trip.Status, dto.TripStatusPublished)
	if err != nil {
		return nil, fmt.Errorf("tripRepository.UpdateStatus: %w", err)
	}
	trip.Status = dto.TripStatusPublished

	return &trip, nil
}
