package auth

import (
	"github.com/go-chi/chi"
)

// Routes related to auth
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", a.Login)
	r.Post("/register", a.Register)
	r.Post("/facebook", a.Facebook)
	r.Post("/google", a.Google)
	r.Post("/email/check", a.CheckEmail)
	r.Post("/username/check", a.CheckUsername)
	r.Get("/verify-email", a.VerifyEmail)
	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.UserCtx)
		// TODO r.Post("/link/google", a.LinkGoogle)
		// TODO r.Post("/link/facebook", a.LinkFacebook)
		r.Post("/logout", a.Logout)
	})
	// r.Get("/sessions", a.Sessions)
	return r
}
