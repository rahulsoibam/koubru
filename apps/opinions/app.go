package opinions

import (
	"database/sql"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/sony/sonyflake"

	"github.com/rahulsoibam/koubru/middleware"

	"github.com/go-redis/redis"
)

// App for user
type App struct {
	DB         *sql.DB
	Cache      *redis.Client
	Middleware *middleware.Middleware
	Flake      *sonyflake.Sonyflake
	Uploader   *s3manager.Uploader
	Sess       *session.Session
}
