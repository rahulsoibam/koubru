package middleware

import (
	"context"
	"net/http"
)

type SearchKeys string

func (m *Middleware) Search(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.FormValue("q")
		q = "%" + q + "%"

		ctx = context.WithValue(ctx, SearchKeys("q"), q)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
