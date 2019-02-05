package categories

import (
	"database/sql"

	"github.com/rahulsoibam/koubru/logger"
	"github.com/rahulsoibam/koubru/middleware"
)

// App for categories
type App struct {
	DB         *sql.DB
	Middleware *middleware.Middleware
	Log        *logger.Logger
}
