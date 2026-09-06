package domain

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (u *TripUsecase) MoveTripPublishedToStarted(
	ctx context.Context,
	tx pgx.Tx,
	req MoveTripPublishToStartRequest,
) (MoveTripPublishToStartResponse, error) {
	tracer := otel.Tracer("TripUsecase")

	ctx, span := tracer.Start(ctx, "TripUsecase.MoveTripPublishedToStarted")
	defer span.End()

	repoReq := toRepositoryMoveTripPublishToStartRequest(req)
	trip, err := u.TripRepo.GetForUpdateByID(ctx, tx, repoReq.ID)
	if err != nil {
		return MoveTripPublishToStartResponse{}, err
	}

	if trip.DriverID != repoReq.DriverID {
		return MoveTripPublishToStartResponse{}, ErrClientNotDriver
	}

	if trip.Status == string(StatusStarted) {
		return MoveTripPublishToStartResponse{
				ID:             trip.ID,
				DriverID:       trip.DriverID,
				FromPoint:      trip.FromPoint,
				ToPoint:        trip.ToPoint,
				DepartureTime:  trip.DepartureTime,
				AvailableSeats: trip.AvailableSeats,
				Status:         trip.Status,
				CreatedAt:      trip.CreatedAt,
			},
			ErrStatusIsStartedAlready
	}

	if trip.Status != string(StatusPublished) {
		return MoveTripPublishToStartResponse{}, ErrNotAllowedCurrentStatusToStarted
	}

	checkResult, err := u.ContractClient.CheckService(ctx, req.ClientID, string(ServiceStarted))
	if err != nil {
		return MoveTripPublishToStartResponse{}, err
	}
	if !checkResult.Allowed {
		return MoveTripPublishToStartResponse{}, errors.New(ErrNotAllowedToStart.Error() + checkResult.Reason)
	}

	tripResp, err := u.TripRepo.MoveTripPublishToStarted(ctx, tx, trip.ID, string(StatusStarted))
	if err != nil {
		return MoveTripPublishToStartResponse{}, err
	}
	domainTripResp := fromRepositoryMoveTripPublishToStartResponse(tripResp)
	//trip.Status = api.TripStatusStarted

	return domainTripResp, nil
}
