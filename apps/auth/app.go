package auth

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/rahulsoibam/koubru/authutils"
	"github.com/rahulsoibam/koubru/middleware"
)

// App for auth
type App struct {
	AuthCache    *redis.Client
	Middleware   *middleware.Middleware
	DB           *sql.DB
	AuthDB       *sql.DB
	Argon2Params *authutils.Params
}
