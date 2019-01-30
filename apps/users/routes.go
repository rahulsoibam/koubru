package users

import (
	"github.com/go-chi/chi"
)

// Routes related to users
func (a App) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Get("/", a.List)
	r.Route("/{username}", func(r chi.Router) {
		r.Get("/", a.Get)
		r.Get("/followers", a.Followers)
		r.Get("/following", a.Following)
		r.Put("/follow", a.Follow)
		r.Delete("/follow", a.Unfollow)
		r.Get("/topics", a.Topics)
		r.Get("/opinions", a.Opinions)
	})
	return r
}
