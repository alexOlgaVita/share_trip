package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (s *TripService) GetTrip(
	ctx context.Context,
	tripId string,
) (GetTripResponse, error) {
	ctx, span := otel.Tracer("TripService").Start(ctx, "TripService.GetTrip")
	defer span.End()

	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*GetTripResponse, error) {
		domainResp, err := s.TripUsecase.GetTrip(ctx, tx, tripId)
		if err != nil {
			return nil, fmt.Errorf("usecase.GetTrip: %w", err)
		}
		serviceResp := fromDomainGetTripResponse(domainResp)

		return &serviceResp, nil
	})

	if err != nil {
		return GetTripResponse{}, fmt.Errorf("failed in transaction: %w", err)
	}

	return *res, nil
}
