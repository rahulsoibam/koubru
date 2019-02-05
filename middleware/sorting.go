package middleware

import (
	"context"
	"net/http"
	"strings"
)

type SortingKeys string

func (m *Middleware) Sorting(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sort := r.FormValue("sort")
		order := r.FormValue("order")

		sort = "time"

		order = strings.ToLower(order)
		switch order {
		case "asc":
		default:
			sort = "desc"
		}

		ctx = context.WithValue(ctx, SortingKeys("sort"), sort)
		ctx = context.WithValue(ctx, SortingKeys("order"), order)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
