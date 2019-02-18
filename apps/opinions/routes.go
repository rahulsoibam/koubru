package opinions

import "github.com/go-chi/chi"

// Routes for opinion
func (a App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.OptionalAuthorization)
		r.Use(a.Middleware.Pagination)
		r.Get("/", a.List)
	})
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.RequireAuthorization)
		r.Post("/", a.Create) // Reply directly to topic
	})
	// with :id
	r.Route("/{opinion_id:[0-9]+}", func(r chi.Router) {
		r.Use(a.Middleware.OpinionCtx)
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.OptionalAuthorization)
			r.Get("/", a.Get)
			r.Get("/followers", a.Followers)
			r.Get("/replies", a.Replies)
			r.Get("/breadcrumbs", a.Breadcrumbs)
		})
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.RequireAuthorization)
			r.Post("/", a.Reply)
			r.Delete("/", a.Delete)
			r.Put("/follow", a.Follow)
			r.Delete("/follow", a.Unfollow)
			r.Post("/report", a.Report)
			r.Post("/vote", a.Vote)
		})
	})
	return r
}
