package domain

import (
	"job4j.ru/share-trip/internal/repository"
)

func fromRepositoryGetTripResponse(response repository.GetTripResponse) GetTripResponse {
	result := GetTripResponse{
		ID:             response.ID,
		DriverID:       response.DriverID,
		FromPoint:      response.FromPoint,
		ToPoint:        response.ToPoint,
		DepartureTime:  response.DepartureTime,
		AvailableSeats: response.AvailableSeats,
		Status:         string(response.Status),
		CreatedAt:      response.CreatedAt,
	}
	return result
}

func toRepositoryCreateTripRequest(req CreateTripRequest) repository.CreateTripRequest {
	return repository.CreateTripRequest{
		ID:             req.ID,
		DriverID:       req.DriverID,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime,
		AvailableSeats: req.AvailableSeats,
		Status:         string(req.Status),
		CreatedAt:      req.CreatedAt,
	}
}

func fromRepositoryCreateTripResponse(resp repository.CreateTripResponse) CreateTripResponse {
	status := TripStatus(resp.Status)
	if status == "" {
		status = StatusDraft
	}
	return CreateTripResponse{
		ID:             resp.ID,
		DriverID:       resp.DriverID,
		FromPoint:      resp.FromPoint,
		ToPoint:        resp.ToPoint,
		DepartureTime:  resp.DepartureTime,
		AvailableSeats: resp.AvailableSeats,
		Status:         status,
		CreatedAt:      resp.CreatedAt,
	}
}

func toRepositoryMoveTripDraftToPublishRequest(request MoveTripDraftToPublishRequest) repository.MoveTripPublishToStartRequest {
	result := repository.MoveTripPublishToStartRequest{
		ID:       request.TripID,
		DriverID: request.ClientID,
	}

	return result
}

func fromRepositoryMoveTripDraftToPublishResponse(response repository.MoveTripDraftToPublishResponse) MoveTripDraftToPublishResponse {
	result := MoveTripDraftToPublishResponse{
		ID:             response.ID,
		DriverID:       response.DriverID,
		FromPoint:      response.FromPoint,
		ToPoint:        response.ToPoint,
		DepartureTime:  response.DepartureTime,
		AvailableSeats: response.AvailableSeats,
		Status:         response.Status,
		CreatedAt:      response.CreatedAt,
	}

	return result
}

func toRepositoryMoveTripPublishToStartRequest(request MoveTripPublishToStartRequest) repository.MoveTripPublishToStartRequest {
	result := repository.MoveTripPublishToStartRequest{
		ID:       request.TripID,
		DriverID: request.ClientID,
	}

	return result
}

func fromRepositoryMoveTripPublishToStartResponse(response repository.MoveTripPublishToStartResponse) MoveTripPublishToStartResponse {
	result := MoveTripPublishToStartResponse{
		ID:             response.ID,
		DriverID:       response.DriverID,
		FromPoint:      response.FromPoint,
		ToPoint:        response.ToPoint,
		DepartureTime:  response.DepartureTime,
		AvailableSeats: response.AvailableSeats,
		Status:         response.Status,
		CreatedAt:      response.CreatedAt,
	}

	return result
}
