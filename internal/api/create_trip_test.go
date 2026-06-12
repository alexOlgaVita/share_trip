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

func TestServer_CreateTrip(t *testing.T) {
	t.Run("success - создание поездки", func(t *testing.T) {
		payload := dto.CreateTripRequest{
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
	})
}
