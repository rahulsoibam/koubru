package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/rahulsoibam/koubru-prod-api/errs"
	"github.com/rahulsoibam/koubru-prod-api/types"

	"github.com/go-chi/chi"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

type TopicKeys string

func (m *Middleware) TopicID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxtopic := types.ContextTopic{}
		topicIDString := chi.URLParam(r, "topic_id")
		topicID, err := strconv.ParseInt(topicIDString, 10, 64)
		if err != nil {
			m.Log.Infoln(err)
			utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
			return
		}
		// err = m.DB.QueryRow("select exists(select 1 from topic where topic_id=$1", topicID).Scan(&exists)
		err = m.DB.QueryRow("select topic_id, title from topic where topic_id = $1", topicID).Scan(&ctxtopic.ID, &ctxtopic.Title)
		if err != nil {
			if err == sql.ErrNoRows {
				m.Log.Infoln(err)
				utils.RespondWithError(w, http.StatusNotFound, errs.TopicNotFound)
			}
			m.Log.Errorln(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}
		ctx = context.WithValue(ctx, TopicKeys("topic_id"), ctxtopic.ID)
		ctx = context.WithValue(ctx, TopicKeys("topic"), ctxtopic)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
