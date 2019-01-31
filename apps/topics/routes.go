package topics

import (
	"github.com/go-chi/chi"
)

// Routes for topics
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", a.List)
	r.Post("/", a.Create)
	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		r.Get("/", a.Get)
		r.Patch("/", a.Patch)
		r.Delete("/", a.Delete)
		r.Get("/followers", a.Followers)
		r.Put("/follow", a.Follow)
		r.Delete("/follow", a.Unfollow)
		r.Post("/report", a.Report)
	})
	return r
}
