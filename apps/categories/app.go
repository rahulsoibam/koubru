package categories

import (
	"database/sql"

	"github.com/rahulsoibam/koubru-prod-api/middleware"
)

// App for categories
type App struct {
	DB         *sql.DB
	Middleware *middleware.Middleware
}
