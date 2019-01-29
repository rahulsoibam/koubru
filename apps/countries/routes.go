package countries

import (
	"github.com/go-chi/chi"
)

// Routes for countries
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", a.List)
	r.Post("/", a.BulkSelect)
	return r
}
