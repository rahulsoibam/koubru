package topics

import (
	"database/sql"

	"github.com/rahulsoibam/koubru/middleware"
)

// App for topics
type App struct {
	DB         *sql.DB
	Middleware *middleware.Middleware
}
