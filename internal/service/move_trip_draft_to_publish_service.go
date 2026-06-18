package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/internal/dto"
	"time"
)

func (s *TripService) MoveTripDraftToPublish(
	ctx context.Context,
	req dto.UpdateTripRequest,
) (*dto.Trip, error) {
	ctx, span := otel.Tracer("TripService").Start(ctx, "TripService.MoveTripDraftToPublish")
	defer span.End()

	started := time.Now()
	result := "success"

	defer func() {
		s.metrics.TripPublishTotal.WithLabelValues(result).Inc()
		s.metrics.TripPublishDuration.WithLabelValues(result).
			Observe(time.Since(started).Seconds())
	}()

	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.Trip, error) {
		resp, err := s.TripUsecase.MoveTripDraftToPublish(ctx, tx, dto.UpdateTripRequest{
			TripID:   req.TripID,
			ClientID: req.ClientID,
		})
		if err != nil {
			if errors.Is(err, domain.ErrTripNotFound) ||
				errors.Is(err, domain.ErrClientNotDriver) ||
				errors.Is(err, domain.ErrNotAllowedCurrentStatusToPublish) ||
				errors.Is(err, domain.ErrStatusIsPublishedAlready) {
				return nil, err
			}
			return nil, fmt.Errorf("usecase.MoveTripDraftToPublish: %w", err)
		}
		return resp, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
