package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"job4j.ru/share-trip/internal/api"
	"job4j.ru/share-trip/internal/api/testauth"
	contractclient "job4j.ru/share-trip/internal/clients/contract"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_MoveTripPublishedToStarted(t *testing.T) {
	t.Parallel()
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт no_active_contract - fail", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{Allowed: false, Reason: string(api.ReasonNoActiveContract)}, nil)

		app := newTestAppWithContractClient(mockClient)

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+string(api.ReasonNoActiveContract))
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт terminated - fail", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{Allowed: false, Reason: string(api.ReasonContractTerminated)}, nil)
		app := newTestAppWithContractClient(mockClient)

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+string(api.ReasonContractTerminated))
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт not started - fail", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{Allowed: false, Reason: string(api.ReasonContractNotStarted)}, nil)
		app := newTestAppWithContractClient(mockClient)

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+string(api.ReasonContractNotStarted))
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт contract_expired - fail", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{Allowed: false, Reason: string(api.ReasonContractExpired)}, nil)
		app := newTestAppWithContractClient(mockClient)

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+string(api.ReasonContractExpired))
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт service_not_in_contract - fail", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{Allowed: false, Reason: string(api.ReasonServiceNotInContract)}, nil)
		app := newTestAppWithContractClient(mockClient)

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+string(api.ReasonServiceNotInContract))
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт service_disabled - fail", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{Allowed: false, Reason: string(api.ReasonServiceDisabled)}, nil)
		app := newTestAppWithContractClient(mockClient)

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+string(api.ReasonServiceDisabled))
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published' - success", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{Allowed: true, Reason: ""}, nil)
		app := newTestAppWithContractClient(mockClient)
		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotRespStart api.MoveTripPublishToStartResponse
		err = json.Unmarshal(respBody, &gotRespStart)
		require.NoError(t, err)
		requireEqualStartededTrip(t, gotRespStart, got)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Started' - success, not content", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{Allowed: true, Reason: ""}, nil)
		app := newTestAppWithContractClient(mockClient)

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotRespStart api.MoveTripPublishToStartResponse
		err = json.Unmarshal(respBody, &gotRespStart)
		require.NoError(t, err)
		requireEqualStartededTrip(t, gotRespStart, got)

		// повторно отправляем на Started
		resp, err = app.Test(req, -1)
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
	t.Run("Перевод поездки в статус 'Started' - from 'Published', Contract Service недоступен - fail", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.started").
			Return(contractclient.CheckResult{}, errors.New("internal server error"))
		app := newTestAppWithContractClient(mockClient)

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
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
	t.Run("Перевод поездки в статус 'Started' - from 'Published', retry - service's avialable after 2 call, retry executes 1 times - success", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()

		// Создаем фейковый внешний API (Contract Service)
		attempts := 0
		externalAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts == 1 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"allowed": true, "reason": ""}`))
		}))

		defer externalAPI.Close()

		cfg := testContractCfg
		cfg.BaseUrl = externalAPI.URL
		app := newTestAppWithContractClient(contractclient.NewClient(cfg))

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, 2, attempts) // Проверяем, что Resty реально сделал 2 запроса
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotRespStart api.MoveTripPublishToStartResponse
		err = json.Unmarshal(respBody, &gotRespStart)
		require.NoError(t, err)
		requireEqualStartededTrip(t, gotRespStart, got)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', retry - service's avialable after 3 call, but retyr executes only 2 times - fail", func(t *testing.T) {
		t.Parallel()
		mockDriverId := uuid.NewString()

		attempts := 0
		externalAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts == 1 || attempts == 2 || attempts == 3 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"allowed": true, "reason": ""}`))
		}))

		defer externalAPI.Close()

		cfg := testContractCfg
		cfg.BaseUrl = externalAPI.URL
		app := newTestAppWithContractClient(contractclient.NewClient(cfg))

		payload := api.CreateTripRequest{
			DriverID:       mockDriverId,
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

		resp, err := app.Test(req, -1)
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

		resp, err = app.Test(req, -1)
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

		payloadStart := api.MoveTripPublishToStartRequest{
			TripID:   tripID,
			ClientID: driverId,
		}

		body, err = json.Marshal(payloadStart)
		require.NoError(t, err)
		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		testauth.WithAuth(req)

		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, 3, attempts) // Проверяем, что Resty реально сделал 2 запроса
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			"internal server error")
	})
}

func requireEqualStartededTrip(t *testing.T, updated api.MoveTripPublishToStartResponse, got api.CreateTripResponse) {
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
		Status:         string(api.StatusStarted),
		CreatedAt:      got.Trip.CreatedAt,
	})
}
