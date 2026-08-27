package configs

import "errors"

type Keycloak struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RequiredRole string
}

func LoadKeycloak() Keycloak {
	return Keycloak{
		Issuer:       Env("KEYCLOAK_ISSUER", "http://localhost:8087/realms/sharetrip"),
		ClientID:     Env("KEYCLOAK_CLIENT_ID", "sharetrip-api"),
		ClientSecret: Env("KEYCLOAK_CLIENT_SECRET", "S8GXiEbeewsHNniYvmUP4AGQM6QbZY3B"),
		RequiredRole: Env("KEYCLOAK_REQUIRED_ROLE", "client"),
	}
}

func (c Keycloak) Validate() error {
	if c.Issuer == "" || c.ClientID == "" || c.ClientSecret == "" {
		return errors.New("KEYCLOAK_ISSUER, KEYCLOAK_CLIENT_ID and KEYCLOAK_CLIENT_SECRET are required")
	}
	if c.RequiredRole == "" {
		return errors.New("KEYCLOAK_REQUIRED_ROLE is required")
	}
	return nil
}
