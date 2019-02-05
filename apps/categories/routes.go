package categories

import (
	"github.com/go-chi/chi"
)

// Routes for categories
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.OptionalAuthorization)
		r.Use(a.Middleware.Pagination)
		r.Use(a.Middleware.Sorting)
		r.Get("/", a.List)
	})
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.RequireAuthorization)
		r.Post("/", a.Create)
		r.Post("/follow", a.BulkFollow)
	})
	r.Route("/{category_id:[0-9]+}", func(r chi.Router) {
		r.Use(a.Middleware.CategoryID)
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.OptionalAuthorization)
			r.Get("/", a.Get)
			r.Group(func(r chi.Router) {
				r.Use(a.Middleware.Pagination)
				r.Get("/topics", a.Topics)
				r.Get("/followers", a.Followers)
			})
		})
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.RequireAuthorization)
			r.Put("/follow", a.Follow)
			r.Delete("/follow", a.Unfollow)
		})
	})
	return r
}
