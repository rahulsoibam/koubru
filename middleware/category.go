package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/rahulsoibam/koubru-prod-api/errs"

	"github.com/go-chi/chi"
	"github.com/rahulsoibam/koubru-prod-api/types"
	"github.com/rahulsoibam/koubru-prod-api/utils"
)

type CategoryKeys string

func (m *Middleware) CategoryID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxcategory := types.ContextCategory{}
		ctx := r.Context()
		categoryIDString := chi.URLParam(r, "category_id")
		categoryID, err := strconv.ParseInt(categoryIDString, 10, 64)
		if err != nil {
			m.Log.Infoln(err)
			utils.RespondWithError(w, http.StatusBadRequest, errs.BadRequest)
			return
		}
		// err = m.DB.QueryRow("select exists(select 1 from category where category_id=$1", categoryID).Scan(&exists)
		err = m.DB.QueryRow("select category_id, name from category where category_id = $1", categoryID).Scan(&ctxcategory.ID, &ctxcategory.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				m.Log.Infoln(err)
				utils.RespondWithError(w, http.StatusNotFound, errs.CategoryNotFound)
			}
			m.Log.Errorln(err)
			utils.RespondWithError(w, http.StatusInternalServerError, errs.InternalServerError)
			return
		}
		ctx = context.WithValue(ctx, CategoryKeys("category_id"), categoryID)
		ctx = context.WithValue(ctx, CategoryKeys("category"), ctxcategory)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
