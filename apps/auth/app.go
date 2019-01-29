package auth

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/rahulsoibam/koubru-prod-api/authutils"
	"github.com/rahulsoibam/koubru-prod-api/middleware"
)

// App for auth
type App struct {
	AuthCache  *redis.Client
	Middleware *middleware.Middleware
	DB         *sql.DB
	// S3Uploader   *s3manager.Uploader
	AuthDB       *sql.DB
	Argon2Params *authutils.Params
}
