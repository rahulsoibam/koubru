package topics

import (
	"github.com/go-chi/chi"
)

// Routes for topics
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.OptionalAuthorization)
		r.Get("/", a.List) // Low priority. The only place this endpoint will be needed is in searching,
		// but searching will be done using ElastiSearch, so no complex queries
	})
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.RequireAuthorization)
		r.Post("/", a.Create)
	})
	r.Route("/{topic_id:[0-9]+}", func(r chi.Router) {
		r.Use(a.Middleware.TopicID)
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.OptionalAuthorization)
			r.Get("/", a.Get)
			r.Get("/followers", a.Followers)
			r.Get("/opinions", a.Opinions)
		})
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.RequireAuthorization)
			r.Patch("/", a.Patch)
			r.Delete("/", a.Delete)
			r.Put("/follow", a.Follow)
			r.Delete("/follow", a.Unfollow)
			r.Post("/report", a.Report)
		})
	})
	return r
}
