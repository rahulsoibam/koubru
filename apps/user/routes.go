package user

import (
	"github.com/go-chi/chi"
)

// Routes for user
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.RequireAuthorization)
		r.Get("/", a.Get) // DONE
		//	r.Patch("/", a.Patch)
		//	r.Delete("/", a.Delete)
		// r.Get("/profile")
		// r.Get|Patch("/profile|settings")
		r.Get("/followers", a.Followers) // DONE
		r.Get("/following", a.Following) // TODO
		r.Get("/opinions", a.Opinions)
		r.Get("/topics", a.Topics) // DONE
	})
	return r
}
