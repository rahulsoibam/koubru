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
type UserCtxKeys int

// Middleware struct for storing redis connection reference
type Middleware struct {
	AuthCache *redis.Client
}

// UserCtx to add user to request context
func (m *Middleware) UserCtx(next http.Handler) http.Handler {
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

		ctx = context.WithValue(ctx, UserCtxKeys(0), userID)
		ctx = context.WithValue(ctx, UserCtxKeys(1), authToken)
		// ctx = context.WithValue(r.Context(), TokenKey, authToken)
		// ctx = context.WithValue(r.Context(), "token", &authToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
