package domain

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"job4j.ru/share-trip/internal/dto"
)

func (u *TripUsecase) GetTrip(
	ctx context.Context,
	tx pgx.Tx,
	tripId string,
) (*dto.Trip, error) {
	tracer := otel.Tracer("TripUsecase")

	ctx, span := tracer.Start(ctx, "TripUsecase.GetTrip")
	defer span.End()

	trip, err := u.TripRepo.GetByID(ctx, tx, tripId)
	if err != nil {
		log.Errorw("s.Repository.Get", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	return trip, nil
}
