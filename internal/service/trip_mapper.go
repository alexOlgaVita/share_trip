package service

import "job4j.ru/share-trip/internal/domain"

func fromDomainGetTripResponse(response domain.GetTripResponse) GetTripResponse {
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

func toDomainCreateTripRequest(request CreateTripRequest) domain.CreateTripRequest {
	result := domain.CreateTripRequest{
		ID:             request.ID,
		DriverID:       request.DriverID,
		FromPoint:      request.FromPoint,
		ToPoint:        request.ToPoint,
		DepartureTime:  request.DepartureTime,
		AvailableSeats: request.AvailableSeats,
	}
	return result
}

func fromDomainCreateTripResponse(response domain.CreateTripResponse) CreateTripResponse {
	result := CreateTripResponse{
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

func toDomainMoveTripDraftToPublishRequest(request MoveTripDraftToPublishRequest) domain.MoveTripDraftToPublishRequest {
	result := domain.MoveTripDraftToPublishRequest{
		TripID:   request.TripID,
		ClientID: request.ClientID,
	}
	return result
}

func fromDomainMoveTripDraftToPublishResponse(response domain.MoveTripDraftToPublishResponse) MoveTripDraftToPublishResponse {
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

func toDomainMoveTripPublishToStartRequest(request MoveTripPublishToStartRequest) domain.MoveTripPublishToStartRequest {
	result := domain.MoveTripPublishToStartRequest{
		TripID:   request.TripID,
		ClientID: request.ClientID,
	}
	return result
}

func fromDomainMoveTripPublishToStartResponse(response domain.MoveTripPublishToStartResponse) MoveTripPublishToStartResponse {
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
