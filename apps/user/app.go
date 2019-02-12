package user

import (
	"database/sql"

	"github.com/rahulsoibam/koubru/middleware"

	"github.com/go-redis/redis"
)

// App for user
type App struct {
	DB         *sql.DB
	Cache      *redis.Client
	Middleware *middleware.Middleware
}
