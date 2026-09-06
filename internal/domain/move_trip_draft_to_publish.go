package domain

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (u *TripUsecase) MoveTripDraftToPublish(
	ctx context.Context,
	tx pgx.Tx,
	req MoveTripDraftToPublishRequest,
) (MoveTripDraftToPublishResponse, error) {
	tracer := otel.Tracer("TripUsecase")

	ctx, span := tracer.Start(ctx, "TripUsecase.MoveTripDraftToPublish")
	defer span.End()

	repoReq := toRepositoryMoveTripDraftToPublishRequest(req)
	trip, err := u.TripRepo.GetForUpdateByID(ctx, tx, repoReq.ID)
	if err != nil {
		return MoveTripDraftToPublishResponse{}, err
	}

	if trip.DriverID != repoReq.DriverID {
		return MoveTripDraftToPublishResponse{}, ErrClientNotDriver
	}

	if trip.Status == string(StatusPublished) {
		return MoveTripDraftToPublishResponse{
				ID:             trip.ID,
				DriverID:       trip.DriverID,
				FromPoint:      trip.FromPoint,
				ToPoint:        trip.ToPoint,
				DepartureTime:  trip.DepartureTime,
				AvailableSeats: trip.AvailableSeats,
				Status:         trip.Status,
				CreatedAt:      trip.CreatedAt,
			},
			ErrStatusIsPublishedAlready
	}

	if trip.Status != string(StatusDraft) {
		return MoveTripDraftToPublishResponse{}, ErrNotAllowedCurrentStatusToPublish
	}

	tripResp, err := u.TripRepo.MoveTripDraftToPublish(ctx, tx, trip.ID, string(StatusPublished))
	if err != nil {
		return MoveTripDraftToPublishResponse{}, err
	}
	domainTripResp := fromRepositoryMoveTripDraftToPublishResponse(tripResp)
	/////trip.Status = string(StatusPublished)

	return domainTripResp, nil
}
