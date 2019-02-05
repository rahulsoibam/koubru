package users

import (
	"database/sql"

	"github.com/rahulsoibam/koubru/logger"

	"github.com/rahulsoibam/koubru/middleware"

	"github.com/go-redis/redis"
)

// App for user
type App struct {
	DB         *sql.DB
	Cache      *redis.Client
	Middleware *middleware.Middleware
	Log        *logger.Logger
}
