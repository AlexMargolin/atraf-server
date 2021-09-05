package middleware

import (
	"context"
	"net/http"
	"strconv"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
)

const (
	LimitParam  = "limit"
	CursorParam = "cursor"
)

const (
	DefaultLimit = 9
	MaxLimit     = 100
)

type PaginationContext struct {
	Limit  int
	Cursor uid.UID
}

type paginationContextKey string

const PaginationContextKey paginationContextKey = "PaginationCtx"

func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pagination := PaginationContext{
			Limit: DefaultLimit,
		}

		limit := r.URL.Query().Get(LimitParam)
		cursor := r.URL.Query().Get(CursorParam)

		if limit != "" {
			if limit, err := strconv.Atoi(limit); err == nil {
				pagination.Limit = limit
			}
		}

		if cursor != "" {
			if cursor, err := uid.FromString(cursor); err == nil {
				pagination.Cursor = cursor
			}
		}

		if pagination.Limit < 1 || pagination.Limit > MaxLimit {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		ctx := context.WithValue(r.Context(), PaginationContextKey, &pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPaginationContext(request *http.Request) *PaginationContext {
	return request.Context().Value(PaginationContextKey).(*PaginationContext)
}
