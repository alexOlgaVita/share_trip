package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/internal/dto"
)

type TripService struct {
	Pool        *pgxpool.Pool
	TripUsecase *domain.TripUsecase
}

func (s *TripService) CreateTrip(
	ctx context.Context,
	req dto.CreateTripRequest,
) (*dto.Trip, error) {
	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.Trip, error) {
		resp, err := s.TripUsecase.CreateTrip(ctx, tx, dto.CreateTripRequest{
			DriverId:       req.DriverId,
			FromPoint:      req.FromPoint,
			ToPoint:        req.ToPoint,
			DepartureTime:  req.DepartureTime,
			AvailableSeats: req.AvailableSeats,
		})
		if err != nil {
			return nil, fmt.Errorf("usecase.CreateTrip: %w", err)
		}
		return resp, nil

	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
