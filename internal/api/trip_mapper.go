package api

import (
	"errors"
	"job4j.ru/share-trip/internal/service"
)

func toServiceCreateTripRequest(request CreateTripRequest) service.CreateTripRequest {
	return service.CreateTripRequest{
		DriverID:       request.DriverID,
		FromPoint:      request.FromPoint,
		ToPoint:        request.ToPoint,
		DepartureTime:  request.DepartureTime,
		AvailableSeats: request.AvailableSeats,
	}
}

// ToCreateTripResponse преобразует ответ сервиса в DTO для API.
// Это явный маппер: он делает семантический слой между service и api.
func fromServiceCreateTripResponse(resp service.CreateTripResponse) (CreateTripResponse, error) {
	// Пример семантической проверки: можно мапить только определённые статусы
	if resp.Status == "" {
		return CreateTripResponse{}, errors.New("trip status is missing")
	}

	return CreateTripResponse{
		Trip: TripResponse{
			ID:             resp.ID,
			DriverID:       resp.DriverID,
			FromPoint:      resp.FromPoint,
			ToPoint:        resp.ToPoint,
			DepartureTime:  resp.DepartureTime,
			AvailableSeats: resp.AvailableSeats,
			Status:         resp.Status,
			CreatedAt:      resp.CreatedAt,
		},
	}, nil
}

func toServiceMoveTripDraftToPublishRequest(request MoveTripDraftToPublishRequest) service.MoveTripDraftToPublishRequest {
	result := service.MoveTripDraftToPublishRequest{
		TripID:   request.TripID,
		ClientID: request.ClientID,
	}

	return result
}

func fromServiceMoveTripDraftToPublishResponse(resp service.MoveTripDraftToPublishResponse) (MoveTripDraftToPublishResponse, error) {
	return MoveTripDraftToPublishResponse{
		Trip: TripResponse{
			ID:             resp.ID,
			DriverID:       resp.DriverID,
			FromPoint:      resp.FromPoint,
			ToPoint:        resp.ToPoint,
			DepartureTime:  resp.DepartureTime,
			AvailableSeats: resp.AvailableSeats,
			Status:         resp.Status,
			CreatedAt:      resp.CreatedAt,
		},
	}, nil
}

func toServiceMoveTripPublishToStartRequest(request MoveTripPublishToStartRequest) service.MoveTripPublishToStartRequest {
	result := service.MoveTripPublishToStartRequest{
		TripID:   request.TripID,
		ClientID: request.ClientID,
	}

	return result
}

func fromServiceMoveTripPublishToStartResponse(resp service.MoveTripPublishToStartResponse) (MoveTripPublishToStartResponse, error) {
	return MoveTripPublishToStartResponse{
		Trip: TripResponse{
			ID:             resp.ID,
			DriverID:       resp.DriverID,
			FromPoint:      resp.FromPoint,
			ToPoint:        resp.ToPoint,
			DepartureTime:  resp.DepartureTime,
			AvailableSeats: resp.AvailableSeats,
			Status:         resp.Status,
			CreatedAt:      resp.CreatedAt,
		},
	}, nil
}
