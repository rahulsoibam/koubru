package users

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/rahulsoibam/koubru-prod-api/middleware"
)

// App for storing database reference
type App struct {
	DB         *sql.DB
	Cache      *redis.Client
	Middleware *middleware.Middleware
}
