package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"job4j.ru/share-trip/internal/domain"
	"strings"
	"time"
)

func (s *TripService) MoveTripPublishToStarted(
	ctx context.Context,
	req MoveTripPublishToStartRequest,
) (MoveTripPublishToStartResponse, error) {
	ctx, span := otel.Tracer("TripService").Start(ctx, "TripService.MoveTripPublishToStarted")
	defer span.End()

	started := time.Now()
	result := "success"

	defer func() {
		s.metrics.TripPublishTotal.WithLabelValues(result).Inc()
		s.metrics.TripPublishDuration.WithLabelValues(result).
			Observe(time.Since(started).Seconds())
	}()

	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*MoveTripPublishToStartResponse, error) {
		domainResp, err := s.TripUsecase.MoveTripPublishedToStarted(ctx, tx, toDomainMoveTripPublishToStartRequest(req))
		if err != nil {
			if errors.Is(err, domain.ErrTripNotFound) ||
				errors.Is(err, domain.ErrClientNotDriver) ||
				errors.Is(err, domain.ErrNotAllowedCurrentStatusToStarted) ||
				errors.Is(err, domain.ErrStatusIsPublishedAlready) ||
				strings.Contains(err.Error(), domain.ErrNotAllowedToStart.Error()) {
				return nil, err
			}
			return nil, fmt.Errorf("usecase.MoveTripPublishToStarted: %w", err)
		}
		serviceResp := fromDomainMoveTripPublishToStartResponse(domainResp)

		// фиксация события в таблице уведомлений в рамках одной транзакции
		err = s.TripUsecase.TripRepo.CreateEvent(ctx, tx, string(TripStarted), req.TripID)
		if err != nil {
			log.Errorw("adding event to outbox after usecase.MoveTripPublishToStarted", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}

		return &serviceResp, nil
	})

	if err != nil {
		return MoveTripPublishToStartResponse{}, fmt.Errorf("failed in transaction: %w", err)
	}

	return *res, nil
}
