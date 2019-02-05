package middleware

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/rahulsoibam/koubru/logger"
)

type Middleware struct {
	AuthCache *redis.Client
	DB        *sql.DB
	Log       *logger.Logger
}
