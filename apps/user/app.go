package user

import (
	"database/sql"

	"github.com/rahulsoibam/koubru-prod-api/logger"

	"github.com/rahulsoibam/koubru-prod-api/middleware"

	"github.com/go-redis/redis"
)

// App for user
type App struct {
	DB         *sql.DB
	Cache      *redis.Client
	Middleware *middleware.Middleware
	Log        *logger.Logger
}
