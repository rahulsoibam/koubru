package middleware

import (
	"database/sql"

	"github.com/go-redis/redis"
)

type Middleware struct {
	AuthCache *redis.Client
	DB        *sql.DB
}
