package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"job4j.ru/share-trip/internal/dto"
)

func (s *TripService) GetTrip(
	ctx context.Context,
	tripId string,
) (*dto.Trip, error) {
	ctx, span := otel.Tracer("TripService").Start(ctx, "TripService.GetTrip")
	defer span.End()

	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.Trip, error) {
		resp, err := s.TripUsecase.GetTrip(ctx, tx, tripId)
		if err != nil {
			return nil, fmt.Errorf("usecase.GetTrip: %w", err)
		}
		return resp, nil

	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
