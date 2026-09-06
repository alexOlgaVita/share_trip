package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	contractclient "job4j.ru/share-trip/internal/clients/contract"
	"job4j.ru/share-trip/internal/observability/logctx"
	"job4j.ru/share-trip/internal/repository"
	"log/slog"
)

type TripUsecase struct {
	TripRepo       *repository.RepoPg
	ContractClient contractclient.Client
}

func (u *TripUsecase) CreateTrip(
	ctx context.Context,
	tx pgx.Tx,
	req CreateTripRequest,
) (CreateTripResponse, error) {
	tracer := otel.Tracer("TripUsecase")

	ctx, span := tracer.Start(ctx, "TripUsecase.CreateTrip")
	defer span.End()

	logger := logctx.Logger(ctx).With(
		slog.String("layer", "usecase"),
		slog.String("usecase", "TripUsecase.CreateTrip"),
		slog.String("client_id", req.DriverID),
	)

	logger.Info("create trip usecase started")

	req.ID = uuid.NewString()
	req.Status = StatusDraft

	repoReq := toRepositoryCreateTripRequest(req)
	tripResp, err := u.TripRepo.Create(ctx, repoReq)
	if err != nil {
		logger.Error(
			"repository create trip failed",
			slog.Any("error", err),
		)
		return CreateTripResponse{}, fmt.Errorf("repoTrip.Create: %w", err)
	}

	domainTripResp := fromRepositoryCreateTripResponse(*tripResp)

	logger.Info(
		"create trip usecase completed",
		slog.String("trip_id", domainTripResp.ID),
	)

	return domainTripResp, nil
}
