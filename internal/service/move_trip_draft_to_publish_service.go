package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"job4j.ru/share-trip/internal/api/dto"
	"job4j.ru/share-trip/internal/domain"
	events "job4j.ru/share-trip/internal/kafka"
	"time"
)

func (s *TripService) MoveTripDraftToPublish(
	ctx context.Context,
	req dto.UpdateTripRequest,
	brokerList []string,
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

		event := events.TripPublished{
			EventID:    uuid.NewString(),
			EventType:  "trip_published",
			TripID:     req.TripID,
			DriverID:   req.DriverId, //пока один и тот же, что и ClientID
			CompanyID:  req.ClientID, //пока один и тот же, что и ClientID
			OccurredAt: time.Now(),   //time.Now().UTC(),
		}
		producer := events.NewProducer(brokerList, "trip.events")
		err = producer.PublishTripPublished(ctx, event)
		if err != nil {
			return nil, fmt.Errorf("usecase.MoveTripDraftToPublish, sending event to kafka error: %w", err)
		}

		defer func() {
			if err := producer.Close(); err != nil {
				log.Errorf("failed to close kafka producer: %v", err)
			}
		}()

		// фиксация события в таблице уведомлений в рамках одной транзакции
		err = s.TripUsecase.TripRepo.CreateEvent(ctx, tx, dto.TripEventPublished, req.TripID)
		if err != nil {
			log.Errorw("adding event to outbox after usecase.MoveTripDraftToPublish", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}

		return resp, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
