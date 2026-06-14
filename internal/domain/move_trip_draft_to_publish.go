package domain

import (
	"context"
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
		return nil, err
	}

	if trip.DriverId != req.ClientID {
		return nil, ErrClientNotDriver
	}

	if trip.Status == dto.TripStatusPublished {
		return trip, ErrStatusIsPublishedAlready
	}

	if trip.Status != dto.TripStatusDraft {
		return nil, ErrNotAllowedCurrentStatusToPublish
	}

	err = u.TripRepo.UpdateStatus(ctx, tx, trip.ID, trip.Status, dto.TripStatusPublished)
	if err != nil {
		return nil, err
	}
	trip.Status = dto.TripStatusPublished

	return trip, nil
}
