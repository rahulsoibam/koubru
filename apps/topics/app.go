package topics

import "database/sql"
import "github.com/rahulsoibam/koubru-prod-api/middleware"

// App for topics
type App struct {
	DB *sql.DB
	m  *middleware.Middleware
}
