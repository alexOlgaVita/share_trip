package domain

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"job4j.ru/share-trip/internal/api/dto"
	contractclient "job4j.ru/share-trip/internal/clients/contract"
)

func (u *TripUsecase) MoveTripPublishedToStarted(
	ctx context.Context,
	tx pgx.Tx,
	req dto.UpdateTripRequest,
) (*dto.Trip, error) {
	tracer := otel.Tracer("TripUsecase")

	ctx, span := tracer.Start(ctx, "TripUsecase.MoveTripPublishedToStarted")
	defer span.End()

	trip, err := u.TripRepo.GetForUpdateByID(ctx, tx, req.TripID)
	if err != nil {
		return nil, err
	}

	if trip.DriverId != req.ClientID {
		return nil, ErrClientNotDriver
	}

	if trip.Status == dto.TripStatusStarted {
		return trip, ErrStatusIsStartedAlready
	}

	if trip.Status != dto.TripStatusPublished {
		return nil, ErrNotAllowedCurrentStatusToStarted
	}

	client := contractclient.ClientFactory()
	checkResult, err := client.CheckService(ctx, req.ClientID, "trip.start")
	if err != nil {
		return nil, err
	}
	if !checkResult.Allowed {
		return nil, errors.New(ErrNotAllowedToStart.Error() + checkResult.Reason)
	}

	err = u.TripRepo.UpdateStatus(ctx, tx, trip.ID, trip.Status, dto.TripStatusStarted)
	if err != nil {
		return nil, err
	}
	trip.Status = dto.TripStatusStarted

	return trip, nil
}
