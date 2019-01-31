package topics

import (
	"github.com/go-chi/chi"
)

// Routes for topics
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(a.Middleware.OptionalUserCtx)
	r.Get("/", a.List)
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.UserCtx)
		r.Post("/", a.Create)
	})
	r.Route("/{topic_id:[0-9]+}", func(r chi.Router) {
		r.Get("/", a.Get)
		r.Get("/followers", a.Followers)
		r.Group(func(r chi.Router) {
			r.Patch("/", a.Patch)
			r.Delete("/", a.Delete)
			r.Put("/follow", a.Follow)
			r.Delete("/follow", a.Unfollow)
			r.Post("/report", a.Report)
		})
	})
	return r
}
