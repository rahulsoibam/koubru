package authutils

import (
	"encoding/base64"
)

// GenerateSecureToken generates a bearer token
func GenerateSecureToken(n uint32) (string, error) {
	// Generate random bytes for token creation
	randomBytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	// Base64 encoded random bytes for token
	bearerToken := base64.RawURLEncoding.EncodeToString(randomBytes)
	return bearerToken, nil
}
