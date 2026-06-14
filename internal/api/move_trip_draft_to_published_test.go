package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"job4j.ru/share-trip/internal/dto"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestServer_MoveTripDraftToPublished_fromDrat_ok(t *testing.T) {
	t.Run("Перевод поездки в статус 'Опубликовано' - from Drat - success", func(t *testing.T) {
		payload := dto.UpdateTripRequest{
			DriverId:       uuid.NewString(),
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)
		require.Equal(t,
			dto.TripRequest{
				DriverId:       payload.DriverId,
				FromPoint:      payload.FromPoint,
				ToPoint:        payload.ToPoint,
				DepartureTime:  payload.DepartureTime,
				AvailableSeats: payload.AvailableSeats,
				Status:         dto.TripStatusDraft,
			},
			dto.TripRequest{
				DriverId:       got.DriverId,
				FromPoint:      got.FromPoint,
				ToPoint:        got.ToPoint,
				DepartureTime:  got.DepartureTime,
				AvailableSeats: got.AvailableSeats,
				Status:         got.Status,
			})

		// при наличии реальной записи, можем проверить позитивный кейс
		tripReq := dto.TripRequest{
			DriverId:       got.DriverId,
			FromPoint:      got.FromPoint,
			ToPoint:        got.ToPoint,
			DepartureTime:  got.DepartureTime,
			AvailableSeats: got.AvailableSeats,
			Status:         dto.TripStatusPublished,
		}

		tripID := got.ID
		driverId := got.DriverId
		payload = dto.UpdateTripRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payload)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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

		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		tripResp := dto.TripRequest{
			DriverId:       tripReq.DriverId,
			FromPoint:      tripReq.FromPoint,
			ToPoint:        tripReq.ToPoint,
			DepartureTime:  tripReq.DepartureTime,
			AvailableSeats: tripReq.AvailableSeats,
			Status:         got.Status,
		}
		require.Equal(t,
			tripReq,
			tripResp)

	})
	t.Run("Перевод поездки в статус 'Опубликовано' - from Published - success", func(t *testing.T) {
		payload := dto.UpdateTripRequest{
			DriverId:       uuid.NewString(),
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)
		require.Equal(t,
			dto.TripRequest{
				DriverId:       payload.DriverId,
				FromPoint:      payload.FromPoint,
				ToPoint:        payload.ToPoint,
				DepartureTime:  payload.DepartureTime,
				AvailableSeats: payload.AvailableSeats,
				Status:         dto.TripStatusDraft,
			},
			dto.TripRequest{
				DriverId:       got.DriverId,
				FromPoint:      got.FromPoint,
				ToPoint:        got.ToPoint,
				DepartureTime:  got.DepartureTime,
				AvailableSeats: got.AvailableSeats,
				Status:         got.Status,
			})

		// при наличии реальной записи, можем проверить позитивный кейс
		tripReq := dto.TripRequest{
			DriverId:       got.DriverId,
			FromPoint:      got.FromPoint,
			ToPoint:        got.ToPoint,
			DepartureTime:  got.DepartureTime,
			AvailableSeats: got.AvailableSeats,
			Status:         dto.TripStatusPublished,
		}

		tripID := got.ID
		driverId := got.DriverId
		payload = dto.UpdateTripRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payload)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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

		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		tripResp := dto.TripRequest{
			DriverId:       tripReq.DriverId,
			FromPoint:      tripReq.FromPoint,
			ToPoint:        tripReq.ToPoint,
			DepartureTime:  tripReq.DepartureTime,
			AvailableSeats: tripReq.AvailableSeats,
			Status:         got.Status,
		}
		require.Equal(t,
			tripReq,
			tripResp)

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
		tripID := ""
		driverId := uuid.NewString()
		payload := dto.UpdateTripRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPut,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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
		tripID := uuid.NewString()
		driverId := ""
		payload := dto.UpdateTripRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPut,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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
		payload := dto.UpdateTripRequest{
			DriverId:       uuid.NewString(),
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)
		require.Equal(t,
			dto.TripRequest{
				DriverId:       payload.DriverId,
				FromPoint:      payload.FromPoint,
				ToPoint:        payload.ToPoint,
				DepartureTime:  payload.DepartureTime,
				AvailableSeats: payload.AvailableSeats,
				Status:         dto.TripStatusDraft,
			},
			dto.TripRequest{
				DriverId:       got.DriverId,
				FromPoint:      got.FromPoint,
				ToPoint:        got.ToPoint,
				DepartureTime:  got.DepartureTime,
				AvailableSeats: got.AvailableSeats,
				Status:         got.Status,
			})

		// попытка публикации несуществующей заявки
		tripID := uuid.NewString()
		driverId := got.DriverId
		payload = dto.UpdateTripRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payload)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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
		payload := dto.UpdateTripRequest{
			DriverId:       uuid.NewString(),
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)
		require.Equal(t,
			dto.TripRequest{
				DriverId:       payload.DriverId,
				FromPoint:      payload.FromPoint,
				ToPoint:        payload.ToPoint,
				DepartureTime:  payload.DepartureTime,
				AvailableSeats: payload.AvailableSeats,
				Status:         dto.TripStatusDraft,
			},
			dto.TripRequest{
				DriverId:       got.DriverId,
				FromPoint:      got.FromPoint,
				ToPoint:        got.ToPoint,
				DepartureTime:  got.DepartureTime,
				AvailableSeats: got.AvailableSeats,
				Status:         got.Status,
			})

		// попытка публикации от имени клиента - не автора заявки
		tripID := got.ID
		driverId := uuid.NewString()
		payload = dto.UpdateTripRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payload)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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
		// TODO - дополнить проверку и добавить кейсы после реализации методов перевода в статусы,
		// из которых недлпустимо осуществлять публикацию
		payload := dto.UpdateTripRequest{
			DriverId:       uuid.NewString(),
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)
		require.Equal(t,
			dto.TripRequest{
				DriverId:       payload.DriverId,
				FromPoint:      payload.FromPoint,
				ToPoint:        payload.ToPoint,
				DepartureTime:  payload.DepartureTime,
				AvailableSeats: payload.AvailableSeats,
				Status:         dto.TripStatusDraft,
			},
			dto.TripRequest{
				DriverId:       got.DriverId,
				FromPoint:      got.FromPoint,
				ToPoint:        got.ToPoint,
				DepartureTime:  got.DepartureTime,
				AvailableSeats: got.AvailableSeats,
				Status:         got.Status,
			})

		//TODO: 2. добавить вызов метода изменения в статус, из которого недопустимо осуществлять публикацию

		// попытка публикации при недопустимом статусе клиента
		tripID := got.ID
		//TODO: 1. заменить на "got.ID" got.DriverId - пока по сути повторяет кейс с неподходящим Клиентом
		driverId := uuid.NewString()
		payload = dto.UpdateTripRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payload)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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
