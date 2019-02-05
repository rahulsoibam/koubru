package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rahulsoibam/koubru-prod-api/errs"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

type UsernameIDKeys string

func (m *Middleware) UsernameID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		username := chi.URLParam(r, "username")

		var usernameID int64
		// err = m.DB.QueryRow("select exists(select 1 from topic where topic_id=$1", topicID).Scan(&exists)
		err := m.DB.QueryRow("select user_id from kuser where username=$1", username).Scan(&usernameID)
		if err != nil {
			if err == sql.ErrNoRows {
				m.Log.Infoln(err)
				utils.RespondWithError(w, http.StatusNotFound, errs.UserNotFound)
			}
			m.Log.Errorln(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}
		ctx = context.WithValue(ctx, UsernameIDKeys("username_id"), usernameID)
		ctx = context.WithValue(ctx, UsernameIDKeys("username"), username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
