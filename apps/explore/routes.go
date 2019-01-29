package explore

import (
	"github.com/go-chi/chi"
)

// Routes for explore page
func (a App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", a.Get)
	return r
}
