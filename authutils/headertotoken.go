package authutils

import (
	"errors"
	"strings"
)

var (
	errAuthHeaderRequired       = errors.New("Authorization header required")
	errAuthBearerSchemeRequired = errors.New("Authorization requires Bearer scheme")
)

// HeaderToTokenString returns token from the Authorization header
func HeaderToTokenString(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errAuthHeaderRequired
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errAuthBearerSchemeRequired
	}
	authToken := authHeader[len("Bearer "):]
	return authToken, nil
}
