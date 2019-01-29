package feed

import (
	"github.com/go-chi/chi"
)

// Routes related to feed
func (a App) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // Some middleware
	r.Get("/", a.Get)
	return r
}
