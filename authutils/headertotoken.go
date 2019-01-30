package authutils

import (
	"errors"
	"strings"
)

var (
	ErrNoHeader = errors.New("Authorization header required")
	ErrNoBearer = errors.New("Authorization requires Bearer scheme")
)

// HeaderToTokenString returns token from the Authorization header
func HeaderToTokenString(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrNoHeader
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrNoBearer
	}
	authToken := authHeader[len("Bearer "):]
	return authToken, nil
}
