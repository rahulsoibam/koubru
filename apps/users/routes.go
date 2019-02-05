package users

import (
	"github.com/go-chi/chi"
)

// Routes for user
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/@{username:^[A-Za-z0-9_.]{3,30}$}", func(r chi.Router) {
		r.Use(a.Middleware.UsernameID)
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.OptionalAuthorization)
			r.Get("/", a.NameGet)                // DONE
			r.Get("/followers", a.NameFollowers) // ERROR
			r.Get("/following", a.NameFollowing) // NEED FIX is_followed NOT WORKING
			r.Get("/topics", a.NameTopics)       // DONE
			//	r.Get("/opinions", a.UsersOpinions)
		})
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.RequireAuthorization)
			r.Put("/follow", a.NameFollow)      // DONE
			r.Delete("/follow", a.NameUnfollow) // DONE
		})
	})
	return r
}
