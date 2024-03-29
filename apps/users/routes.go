package users

import (
	"github.com/go-chi/chi"
)

// Routes for user
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.Pagination)
		r.Use(a.Middleware.Search)
		r.Get("/", a.List)
	})
	r.Route("/@{username:^[A-Za-z0-9_.]{3,30}$}", func(r chi.Router) {
		r.Use(a.Middleware.UsernameID)
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.RequireAuthorization)
			r.Put("/follow", a.Follow)      // DONE
			r.Delete("/follow", a.Unfollow) // DONE
		})
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.OptionalAuthorization)
			r.Get("/", a.Get) // DONE
		})
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.OptionalAuthorization)
			r.Use(a.Middleware.Pagination)
			r.Get("/followers", a.Followers)
			r.Get("/following", a.Following)
			r.Get("/opinions", a.Opinions)
			r.Get("/topics", a.Topics)
		})
	})
	return r
}
