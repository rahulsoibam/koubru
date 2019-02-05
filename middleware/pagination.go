package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type PaginationKeys string

func (m *Middleware) Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		perPage := r.FormValue("per_page")
		page := r.FormValue("page")

		perPageInt, err := strconv.Atoi(perPage)
		if err != nil || perPageInt <= 0 || perPageInt >= 150 {
			perPageInt = 30
		}
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt <= 0 {
			pageInt = 1
		}

		dbOffset := (pageInt - 1) * perPageInt

		ctx = context.WithValue(ctx, PaginationKeys("per_page"), perPageInt)
		ctx = context.WithValue(ctx, PaginationKeys("page"), pageInt)
		ctx = context.WithValue(ctx, PaginationKeys("db_offset"), dbOffset)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
