package opinions

import "github.com/go-chi/chi"

// Routes for opinion
func (a App) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // Some middleware
	r.Get("/", a.List)
	r.Post("/", a.Create)
	// with :id
	r.Route("/{opinion_id}", func(r chi.Router) {
		r.Get("/", a.Get)
		r.Delete("/", a.Delete)
		r.Get("/followers", a.Followers)
		r.Put("/follow", a.Follow)
		r.Delete("/follow", a.Unfollow)
		r.Get("/replies", a.Replies)
		r.Post("/report", a.Report)
		r.Put("/vote", a.Vote)
	})
	return r
}
