package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/rahulsoibam/koubru-prod-api/authutils"

	"github.com/go-redis/redis"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// UserCtxKeys stores key to use when accessing context values
type AuthKeys string

type Middleware struct {
	AuthCache *redis.Client
}

func (m *Middleware) RequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var userID int64
		var err error
		authHeader := r.Header.Get("Authorization")
		authToken, err := authutils.HeaderToTokenString(authHeader)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Get user id from redis session store cache
		response, err := m.AuthCache.Get(authToken).Result()
		if err == redis.Nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to perfrom the following action. Please login or signup")
			return
		} else if err != nil {
			// If there is an error fetching from the cache, return an internal server error
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Convert to integer
		userID, _ = strconv.ParseInt(response, 10, 64)

		ctx = context.WithValue(ctx, AuthKeys("userID"), userID)
		ctx = context.WithValue(ctx, AuthKeys("authToken"), authToken)
		// ctx = context.WithValue(r.Context(), TokenKey, authToken)
		// ctx = context.WithValue(r.Context(), "token", &authToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) OptionalAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var userID int64
		var err error
		authHeader := r.Header.Get("Authorization")
		authToken, err := authutils.HeaderToTokenString(authHeader)
		if err != nil {
			if err == authutils.ErrNoHeader {
				ctx = context.WithValue(ctx, AuthKeys("userID"), 0)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Get user id from redis session store cache
		response, err := m.AuthCache.Get(authToken).Result()
		if err == redis.Nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to perfrom the following action. Please login or signup")
			return
		} else if err != nil {
			// If there is an error fetching from the cache, return an internal server error
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Convert to integer
		userID, _ = strconv.ParseInt(response, 10, 64)

		ctx = context.WithValue(ctx, AuthKeys("userID"), userID)
		ctx = context.WithValue(ctx, AuthKeys("authToken"), authToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
