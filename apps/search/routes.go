package search

import (
	"github.com/go-chi/chi"
)

// Routes for search
func (a App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/users", a.Users)
	r.Get("/topics", a.Topics)
	r.Get("/categories", a.Categories)
	return r
}
