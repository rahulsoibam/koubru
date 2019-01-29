package authutils

import (
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/go-redis/redis"
)

// Token struct for storing token
type Token struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
}

// Authenticate generates a token, stores the session, and returns the token
func Authenticate(authCache *redis.Client, authDB *sql.DB, userID int64, userAgent string) (interface{}, error) {
	// Generate random bytes for token creation
	randomBytes, err := GenerateRandomBytes(256)
	if err != nil {
		return nil, err
	}
	// Base64 encoded random bytes for token
	bearerToken := base64.RawURLEncoding.EncodeToString(randomBytes)
	// Expiry time of token
	expiry := 60 * 60 * 24 * 30 * time.Second

	// Store token as session in Redis
	err = authCache.Set(bearerToken, userID, expiry).Err()
	if err != nil {
		return nil, err
	}

	_, err = authDB.Exec("INSERT INTO Session (user_id, token, user_agent) VALUES ($1, $2, $3)", userID, bearerToken, userAgent)
	if err != nil {
		return nil, err
	}

	token := Token{
		TokenType:   "Bearer",
		AccessToken: bearerToken,
		Expires:     expiry.Nanoseconds() / 1e9,
	}

	return &token, nil

}
