package topics

import (
	"database/sql"

	"github.com/rahulsoibam/koubru-prod-api/logger"
	"github.com/rahulsoibam/koubru-prod-api/middleware"
)

// App for topics
type App struct {
	DB         *sql.DB
	Middleware *middleware.Middleware
	Log        *logger.Logger
}
