package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"job4j.ru/share-trip/internal/dto"
)

func (s *TripService) MoveTripDraftToPublish(
	ctx context.Context,
	req dto.UpdateTripRequest,
) (*dto.Trip, error) {
	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.Trip, error) {
		resp, err := s.TripUsecase.MoveTripDraftToPublish(ctx, tx, dto.UpdateTripRequest{
			TripID:   req.TripID,
			ClientID: req.ClientID,
		})
		if err != nil {
			return nil, fmt.Errorf("usecase.MoveTripDraftToPublish: %w", err)
		}
		return resp, nil

	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
