package categories

import (
	"github.com/go-chi/chi"
)

// Routes for categories
func (a App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", a.List)
	r.Post("/", a.Create)
	r.Post("/follow", a.BulkFollow)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", a.Get)
		r.Put("/follow", a.Follow)
		r.Delete("/follow", a.Unfollow)
	})
	return r
}
