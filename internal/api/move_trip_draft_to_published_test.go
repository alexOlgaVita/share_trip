package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"job4j.ru/share-trip/internal/api"
	"job4j.ru/share-trip/internal/api/testauth"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestServer_MoveTripDraftToPublished_fromDrat_ok(t *testing.T) {
	t.Parallel()
	t.Run("Перевод поездки в статус 'Опубликовано' - from Drat - success", func(t *testing.T) {
		t.Parallel()
		payload := api.CreateTripRequest{
			DriverID:       uuid.NewString(),
			FromPoint:      "Дубаи",
			ToPoint:        "Екатеринбург",
			DepartureTime:  "2027-01-02 15:04:00",
			AvailableSeats: "1",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got api.CreateTripResponse
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

		tripID := got.Trip.ID
		driverId := got.Trip.DriverID
		payloadPublish := api.MoveTripDraftToPublishRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadPublish)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/publish",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotResp api.MoveTripDraftToPublishResponse
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, gotResp, got)
	})
	t.Run("Перевод поездки в статус 'Опубликовано' - from Published - success", func(t *testing.T) {
		t.Parallel()
		payload := api.CreateTripRequest{
			DriverID:       uuid.NewString(),
			FromPoint:      "Дубаи",
			ToPoint:        "Екатеринбург",
			DepartureTime:  "2027-01-02 15:04:00",
			AvailableSeats: "1",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got api.CreateTripResponse
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

		tripID := got.Trip.ID
		driverId := got.Trip.DriverID
		payloadPublish := api.MoveTripDraftToPublishRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadPublish)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/publish",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotResp api.MoveTripDraftToPublishResponse
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, gotResp, got)

		// повторно отправляем на публикацию
		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t,
			string(respBody),
			"")
	})
	t.Run("Перевод поездки в статус 'Опубликовано' - tripId empty - fail", func(t *testing.T) {
		t.Parallel()
		tripID := ""
		driverId := uuid.NewString()
		payload := api.MoveTripDraftToPublishRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPut,
			"/trip/publish",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t,
			string(respBody),
			"tripID is required")

	})
	t.Run("Перевод поездки в статус 'Опубликовано' - clientId empty - fail", func(t *testing.T) {
		t.Parallel()
		tripID := uuid.NewString()
		driverId := ""
		payload := api.MoveTripDraftToPublishRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPut,
			"/trip/publish",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t,
			string(respBody),
			"clientID is required")

	})
	t.Run("Перевод поездки в статус 'Опубликовано' - trip, having tripId, doesn't exist - fail", func(t *testing.T) {
		t.Parallel()
		payload := api.CreateTripRequest{
			DriverID:       uuid.NewString(),
			FromPoint:      "Дубаи",
			ToPoint:        "Екатеринбург",
			DepartureTime:  "2027-01-02 15:04:00",
			AvailableSeats: "1",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got api.CreateTripResponse
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

		// попытка публикации несуществующей заявки
		tripID := uuid.NewString()
		driverId := got.Trip.DriverID
		payloadPublish := api.MoveTripDraftToPublishRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadPublish)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/publish",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			"internal server error")
	})
	t.Run("Перевод поездки в статус 'Опубликовано' - clientId isn't equal driverId - fail", func(t *testing.T) {
		t.Parallel()
		payload := api.CreateTripRequest{
			DriverID:       uuid.NewString(),
			FromPoint:      "Дубаи",
			ToPoint:        "Екатеринбург",
			DepartureTime:  "2027-01-02 15:04:00",
			AvailableSeats: "1",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got api.CreateTripResponse
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

		tripID := got.Trip.ID
		driverId := uuid.NewString()
		payloadPublish := api.MoveTripDraftToPublishRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadPublish)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/publish",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			"client is not driver of this trip")
	})
	t.Run("Перевод поездки в статус 'Опубликовано' - trip status isn't equal draft or published  - fail", func(t *testing.T) {
		t.Parallel()
		// TODO - дополнить проверку и добавить кейсы после реализации методов перевода в статусы,
		// из которых недлпустимо осуществлять публикацию
		payload := api.CreateTripRequest{
			DriverID:       uuid.NewString(),
			FromPoint:      "Дубаи",
			ToPoint:        "Екатеринбург",
			DepartureTime:  "2027-01-02 15:04:00",
			AvailableSeats: "1",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req) // <-- обязательно: Без этого middleware вернёт 401 даже с mock-сервером.

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got api.CreateTripResponse
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

		//TODO: 2. добавить вызов метода изменения в статус, из которого недопустимо осуществлять публикацию

		// попытка публикации при недопустимом статусе клиента
		tripID := got.Trip.ID
		//TODO: 1. заменить на "got.ID" got.DriverId - пока по сути повторяет кейс с неподходящим Клиентом
		driverId := uuid.NewString()
		payloadPublish := api.MoveTripDraftToPublishRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadPublish)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/publish",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		//TODO: 3. Заменить "http.StatusForbidden" на "http.StatusConflict"
		require.Equal(t, http.StatusForbidden, resp.StatusCode)
		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		//TODO: 4. Заменить "client is not driver of this trip" на "current status is not allowed for publish"
		require.Equal(t,
			string(respBody),
			"client is not driver of this trip")
	})
}

func requireEqualCreatedDraftTrip(t *testing.T, got api.CreateTripResponse, payload api.CreateTripRequest) {
	t.Helper()
	require.NotEmpty(t, got.Trip.ID)
	require.Equal(t, api.TripResponse{
		ID:             got.Trip.ID,
		DriverID:       payload.DriverID,
		FromPoint:      payload.FromPoint,
		ToPoint:        payload.ToPoint,
		DepartureTime:  payload.DepartureTime,
		AvailableSeats: payload.AvailableSeats,
		Status:         string(api.StatusDraft),
		CreatedAt:      got.Trip.CreatedAt,
	}, got.Trip)
}

func requireEqualPublishedTrip(t *testing.T, updated api.MoveTripDraftToPublishResponse, got api.CreateTripResponse) {
	t.Helper()
	require.NotEmpty(t, updated.Trip.ID)
	require.NotEmpty(t, updated.Trip.CreatedAt)

	dateTimeFormat, err := dateToUserFormat(got.Trip.DepartureTime)
	require.NoError(t, err)

	require.Equal(t, api.TripResponse{
		ID:             updated.Trip.ID,
		DriverID:       updated.Trip.DriverID,
		FromPoint:      updated.Trip.FromPoint,
		ToPoint:        updated.Trip.ToPoint,
		DepartureTime:  updated.Trip.DepartureTime,
		AvailableSeats: updated.Trip.AvailableSeats,
		Status:         updated.Trip.Status,
		CreatedAt:      got.Trip.CreatedAt,
	}, api.TripResponse{
		ID:             got.Trip.ID,
		DriverID:       got.Trip.DriverID,
		FromPoint:      got.Trip.FromPoint,
		ToPoint:        got.Trip.ToPoint,
		DepartureTime:  dateTimeFormat,
		AvailableSeats: got.Trip.AvailableSeats,
		Status:         string(api.StatusPublished),
		CreatedAt:      got.Trip.CreatedAt,
	})
}

func dateToUserFormat(dateTimeIn string) (string, error) {
	templateDate := "2006-01-02 15:04:05"
	dateTimeOutParse, err := time.Parse(templateDate, dateTimeIn)
	if err != nil {
		return "", err
	}
	outputLayout := "01-02-2006 15:04"
	dateTimeOut := dateTimeOutParse.Format(outputLayout)
	return dateTimeOut, nil
}
