package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/rahulsoibam/koubru/authutils"
	"github.com/rahulsoibam/koubru/errs"

	"github.com/go-redis/redis"
	"github.com/rahulsoibam/koubru/utils"
)

type AuthKeys string

func (m *Middleware) RequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var userID int64
		var err error
		authHeader := r.Header.Get("Authorization")
		authToken, err := authutils.HeaderToTokenString(authHeader)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusUnauthorized, err) // Directly returning err to user is harmless here. Custom token with harmless message
			return
		}

		// Get user id from redis session store cache
		response, err := m.AuthCache.Get(authToken).Result()
		if err == redis.Nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
			return
		} else if err != nil {
			log.Println(err)
			// If there is an error fetching from the cache, return an internal server error
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}

		// Convert to integer
		userID, _ = strconv.ParseInt(response, 10, 64)

		ctx = context.WithValue(ctx, AuthKeys("user_id"), userID)
		ctx = context.WithValue(ctx, AuthKeys("auth_token"), authToken)
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
			if err == errs.NoHeader {
				ctx = context.WithValue(ctx, AuthKeys("user_id"), 0)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			log.Println(err)
			utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
			return
		}

		// Get user id from redis session store cache
		response, err := m.AuthCache.Get(authToken).Result()
		if err == redis.Nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusUnauthorized, errs.Unauthorized)
			return
		} else if err != nil {
			// If there is an error fetching from the cache, return an internal server error
			log.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}

		// Convert to integer
		userID, err = strconv.ParseInt(response, 10, 64)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusBadRequest, errs.Unauthorized)
		}

		ctx = context.WithValue(ctx, AuthKeys("user_id"), userID)
		ctx = context.WithValue(ctx, AuthKeys("auth_token"), authToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
