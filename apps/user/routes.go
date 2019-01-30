package user

import (
	"github.com/go-chi/chi"
)

// Routes for user
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
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
	return r
}
