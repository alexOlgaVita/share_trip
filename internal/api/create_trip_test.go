package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"job4j.ru/share-trip/internal/api"
	"job4j.ru/share-trip/internal/api/testauth"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestServer_CreateTrip(t *testing.T) {
	t.Parallel()
	t.Run("success - создание поездки", func(t *testing.T) {
		t.Parallel()
		payload := api.TripRequest{
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
		require.NotEmpty(t, got.Trip.ID)

		expected := api.TripResponse{
			ID:             got.Trip.ID,
			DriverID:       payload.DriverID,
			FromPoint:      payload.FromPoint,
			ToPoint:        payload.ToPoint,
			DepartureTime:  payload.DepartureTime,
			AvailableSeats: payload.AvailableSeats,
			Status:         string(api.StatusDraft),
			CreatedAt:      got.Trip.CreatedAt,
		}
		require.Equal(t, expected, got.Trip)
	})
}
