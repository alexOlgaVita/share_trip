package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"job4j.ru/share-trip/configs"
	"job4j.ru/share-trip/internal/api/dto"
	"job4j.ru/share-trip/internal/api/testauth"
	contractclient "job4j.ru/share-trip/internal/clients/contract"
	"job4j.ru/share-trip/internal/domain"
	"job4j.ru/share-trip/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_MoveTripPublishedToStarted(t *testing.T) {
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт no_active_contract - fail", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{Allowed: false, Reason: dto.ReasonNoActiveContract}, nil)

		// Подменяем глобальную фабрику на нашу мок-версию
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+dto.ReasonNoActiveContract)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт terminated - fail", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{Allowed: false, Reason: dto.ReasonContractTerminated}, nil)
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+dto.ReasonContractTerminated)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт not started - fail", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{Allowed: false, Reason: dto.ReasonContractNotStarted}, nil)
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+dto.ReasonContractNotStarted)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт contract_expired - fail", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{Allowed: false, Reason: dto.ReasonContractExpired}, nil)
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+dto.ReasonContractExpired)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт service_not_in_contract - fail", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{Allowed: false, Reason: dto.ReasonServiceNotInContract}, nil)
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+dto.ReasonServiceNotInContract)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published', контракт service_disabled - fail", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{Allowed: false, Reason: dto.ReasonServiceDisabled}, nil)
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		require.Equal(t, http.StatusConflict, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			domain.ErrNotAllowedToStart.Error()+dto.ReasonServiceDisabled)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Published' - success", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{Allowed: true, Reason: ""}, nil)
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}
		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualStartededTrip(t, got, gotResp)
	})
	t.Run("Перевод поездки в статус 'Started' - from 'Started' - success, not content", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{Allowed: true, Reason: ""}, nil)
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualStartededTrip(t, got, gotResp)

		// повторно отправляем на Started
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
	t.Run("Перевод поездки в статус 'Started' - from 'Published', Contract Service недоступен - fail", func(t *testing.T) {
		mockDriverId := uuid.NewString()
		mockClient := new(mocks.Client)
		mockClient.On("CheckService", mock.Anything, mockDriverId, "trip.start").
			Return(contractclient.CheckResult{}, errors.New("internal server error"))
		oldFactory := contractclient.ClientFactory
		defer func() { contractclient.ClientFactory = oldFactory }() // Возвращаем всё как было
		contractclient.ClientFactory = func() contractclient.Client {
			return mockClient
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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
	t.Run("Перевод поездки в статус 'Started' - from 'Published', retry - service's avialable after 2 call, retry executes 1 times - success", func(t *testing.T) {
		mockDriverId := uuid.NewString()

		// Создаем фейковый внешний API (Contract Service)
		attempts := 0
		externalAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts == 1 {
				/////
				//w.Header().Set("Content-Type", "application/json")
				/////
				w.WriteHeader(http.StatusServiceUnavailable)
				/////
				//_, _ = w.Write([]byte(`{"allowed": false, "reason": "service unavailable"}`)) // ОБЯЗАТЕЛЬНО пишем тело
				/////
				return
			}
			/////
			//fmt.Println("MOCK: Request received!")
			/////
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"allowed": true, "reason": ""}`))
		}))

		defer externalAPI.Close()
		// Сохраняем старую, чтобы вернуть
		oldFactory := contractclient.ClientFactory
		t.Cleanup(func() {
			contractclient.ClientFactory = oldFactory
		})

		contractclient.ClientFactory = func() contractclient.Client {
			cfgContract, err := configs.LoadContract()
			if err != nil {
				t.Fatalf("Критическая ошибка: %v", err)
			}
			if err := cfgContract.Validate(); err != nil {
				t.Fatalf("Критическая ошибка: %v", err)
			}
			// Переопределяем только то, что нужно для теста
			cfgContract.BaseUrl = externalAPI.URL
			return contractclient.NewClient(cfgContract)
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		require.Equal(t, 2, attempts) // Проверяем, что Resty реально сделал 2 запроса
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualStartededTrip(t, got, gotResp)

	})

	t.Run("Перевод поездки в статус 'Started' - from 'Published', retry - service's avialable after 3 call, but retyr executes only 2 times - fail", func(t *testing.T) {
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

		oldFactory := contractclient.ClientFactory // Сохраняем старую, чтобы вернуть
		t.Cleanup(func() {
			contractclient.ClientFactory = oldFactory
		})

		contractclient.ClientFactory = func() contractclient.Client {
			cfgContract, err := configs.LoadContract()
			if err != nil {
				t.Fatalf("Критическая ошибка: %v", err)
			}
			if err := cfgContract.Validate(); err != nil {
				t.Fatalf("Критическая ошибка: %v", err)
			}
			// Переопределяем только то, что нужно для теста
			cfgContract.BaseUrl = externalAPI.URL
			return contractclient.NewClient(cfgContract)
		}

		payload := dto.UpdateTripRequest{
			DriverId:       mockDriverId,
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

		var got dto.Trip
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftTrip(t, got, payload)

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

		var gotResp dto.Trip
		err = json.Unmarshal(respBody, &gotResp)
		require.NoError(t, err)
		requireEqualPublishedTrip(t, got, gotResp)

		// перевод поездки в статус "started"
		req, err = http.NewRequest(
			http.MethodPut,
			"/trip/started",
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

		require.Equal(t, 3, attempts) // Проверяем, что Resty реально сделал 2 запроса
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t,
			string(respBody),
			"internal server error")
	})
}

func requireEqualStartededTrip(t *testing.T, created dto.Trip, got dto.Trip) {
	t.Helper()
	require.NotEmpty(t, got.ID)
	require.NotEmpty(t, got.CreatedAt)

	dateTimeFormat, err := dateToUserFormat(created.DepartureTime)
	require.NoError(t, err)

	require.Equal(t, dto.Trip{
		ID:             created.ID,
		DriverId:       created.DriverId,
		FromPoint:      created.FromPoint,
		ToPoint:        created.ToPoint,
		DepartureTime:  dateTimeFormat,
		AvailableSeats: created.AvailableSeats,
		Status:         dto.TripStatusStarted,
		CreatedAt:      got.CreatedAt,
	}, got)
}
