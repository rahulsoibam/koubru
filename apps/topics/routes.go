package topics

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

// Routes for topics
func (a *App) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(a.m.UserCtx)
	r.Get("/", a.List)
	r.Post("/", a.Create)
	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		r.Use(a.TopicCtx)
		r.Get("/", a.Get)
		r.Patch("/", a.Patch)
		r.Delete("/", a.Delete)
		r.Get("/followers", a.Followers)
		r.Put("/follow", a.Follow)
		r.Delete("/follow", a.Unfollow)
		r.Post("/report", a.Report)
	})
	return r
}

// TopicCtx passes the current article to subsequent requests
func (a *App) TopicCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var topicID int64
		if topicIDString := chi.URLParam(r, "id"); topicIDString != "" {
			var err error
			topicID, err = strconv.ParseInt(topicIDString, 10, 64)
			if err != nil {
				utils.RespondWithError(w, http.StatusNotFound, "Invalid Topic URL")
				return
			}
		}
		ctx := context.WithValue(r.Context(), "topic_id", topicID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
