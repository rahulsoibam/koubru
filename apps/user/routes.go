package user

import (
	"github.com/go-chi/chi"
)

// Routes for user
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.UserCtx)
		r.Get("/", a.Get) // DONE
		//	r.Patch("/", a.Patch)
		//	r.Delete("/", a.Delete)
		// r.Get("/profile")
		// r.Get|Patch("/profile|settings")
		r.Get("/followers", a.Followers) // DONE
		r.Get("/following", a.Following) // TODO
		//	r.Get("/opinions", a.Opinions)
		r.Get("/topics", a.Topics) // DONE
	})
	r.Route("/@{username:^[A-Za-z0-9_.]{3,30}$}", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.OptionalUserCtx)
			r.Get("/", a.NameGet)                // DONE
			r.Get("/followers", a.NameFollowers) // ERROR
			r.Get("/following", a.NameFollowing) // NEED FIX is_followed NOT WORKING
			r.Get("/topics", a.NameTopics)       // DONE
			//	r.Get("/opinions", a.UsersOpinions)
		})
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.UserCtx)
			r.Put("/follow", a.NameFollow)      // DONE
			r.Delete("/follow", a.NameUnfollow) // DONE
		})
	})
	return r
}
