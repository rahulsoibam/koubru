package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/rahulsoibam/koubru/errs"

	"github.com/go-chi/chi"
	"github.com/rahulsoibam/koubru/types"
	"github.com/rahulsoibam/koubru/utils"
)

type OpinionKeys string

func (m *Middleware) OpinionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxOpinion := types.ContextOpinion{}
		ctx := r.Context()
		opinionIDString := chi.URLParam(r, "opinion_id")
		opinionID, err := strconv.ParseInt(opinionIDString, 10, 64)
		if err != nil {
			m.Log.Infoln(err)
			utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
			return
		}
		// err = m.DB.QueryRow("select exists(select 1 from category where category_id=$1", categoryID).Scan(&exists)
		err = m.DB.QueryRow("select opinion_id, topic_id, creator_id from topic where topic_id = $1", opinionID).Scan(&ctxOpinion.ID, &ctxOpinion.TopicID, &ctxOpinion.CreatorID)
		if err != nil {
			if err == sql.ErrNoRows {
				m.Log.Infoln(err)
				utils.RespondWithError(w, http.StatusNotFound, errs.OpinionNotFound)
				return
			}
			m.Log.Errorln(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}
		ctx = context.WithValue(ctx, OpinionKeys("opinion_id"), opinionID)
		ctx = context.WithValue(ctx, OpinionKeys("ctx_opinion"), ctxOpinion)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
