package user

import (
	"github.com/go-chi/chi"
)

// Routes for user
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.UserCtx)
		r.Get("/", a.Get)
		r.Patch("/", a.Patch)
		r.Delete("/", a.Delete)
		// r.Get("/profile")
		// r.Get|Patch("/profile|settings")
		r.Get("/followers", a.Followers)
		r.Get("/following", a.Following)
		r.Get("/opinions", a.Opinions)
		r.Get("/topics", a.Topics)
	})
	r.Route("s/{username}", func(r chi.Router) {
		r.Get("/", a.UsersGet)
		r.Get("/followers", a.UsersFollowers)
		r.Get("/following", a.UsersFollowing)
		r.Get("/topics", a.UsersTopics)
		r.Get("/opinions", a.UsersOpinions)
		r.Group(func(r chi.Router) {
			r.Use(a.Middleware.UserCtx)
			r.Put("/follow", a.FollowUser)
			r.Delete("/follow", a.UnfollowUser)
		})
	})
	return r
}
