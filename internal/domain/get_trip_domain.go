package domain

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (u *TripUsecase) GetTrip(
	ctx context.Context,
	tx pgx.Tx,
	tripId string,
) (GetTripResponse, error) {
	tracer := otel.Tracer("TripUsecase")

	ctx, span := tracer.Start(ctx, "TripUsecase.GetTrip")
	defer span.End()

	trip, err := u.TripRepo.GetByID(ctx, tx, tripId)
	if err != nil {
		log.Errorw("s.Repository.Get", err)
		return GetTripResponse{}, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return fromRepositoryGetTripResponse(*trip), nil
}
