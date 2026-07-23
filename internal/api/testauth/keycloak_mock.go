package testauth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"job4j.ru/share-trip/internal/middleware"
)

func NewKeycloakMock(clientID string, roles ...string) (*httptest.Server, middleware.KeycloakConfig) {
	accessToken := mustAccessToken(clientID, roles...)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost ||
			!strings.HasSuffix(r.URL.Path, "/protocol/openid-connect/token") {
			http.NotFound(w, r)
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]string{
			"access_token": accessToken,
		})
	}))

	cfg := middleware.KeycloakConfig{
		Issuer:     srv.URL + "/realms/sharetrip",
		ClientID:   clientID,
		HTTPClient: srv.Client(),
	}

	return srv, cfg
}

func mustAccessToken(clientID string, roles ...string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload, _ := json.Marshal(map[string]any{
		"sub": "test-user-id",
		"resource_access": map[string]any{
			clientID: map[string]any{"roles": roles},
		},
	})
	body := base64.RawURLEncoding.EncodeToString(payload)
	return header + "." + body + ".test-signature"
}

const RefreshToken = "test-refresh-token"

func WithAuth(req *http.Request) {
	req.Header.Set(middleware.RefreshTokenHeader, RefreshToken)
}
