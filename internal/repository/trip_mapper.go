package repository

//func toTripResponse(resp Trip) Trip {
//	return Trip{
//		ID:             resp.ID,
//		DriverID:       resp.DriverID,
//		FromPoint:      resp.FromPoint,
//		ToPoint:        resp.ToPoint,
//		DepartureTime:  resp.DepartureTime,
//		AvailableSeats: resp.AvailableSeats,
//		Status:         resp.Status,
//		CreatedAt:      resp.CreatedAt,
//	}
//}

func toMoveTripDraftToPublishResponse(resp UpdateTripResponse) MoveTripDraftToPublishResponse {
	return MoveTripDraftToPublishResponse(resp)
	//return MoveTripDraftToPublishResponse{
	//	ID:             resp.ID,
	//	DriverID:       resp.DriverID,
	//	FromPoint:      resp.FromPoint,
	//	ToPoint:        resp.ToPoint,
	//	DepartureTime:  resp.DepartureTime,
	//	AvailableSeats: resp.AvailableSeats,
	//	Status:         resp.Status,
	//	CreatedAt:      resp.CreatedAt,
	//}
}

func toMoveTripPublishToStartResponse(resp UpdateTripResponse) MoveTripPublishToStartResponse {
	return MoveTripPublishToStartResponse(resp)
	//return MoveTripPublishToStartResponse{
	//	ID:             resp.ID,
	//	DriverID:       resp.DriverID,
	//	FromPoint:      resp.FromPoint,
	//	ToPoint:        resp.ToPoint,
	//	DepartureTime:  resp.DepartureTime,
	//	AvailableSeats: resp.AvailableSeats,
	//	Status:         resp.Status,
	//	CreatedAt:      resp.CreatedAt,
	//}
}
