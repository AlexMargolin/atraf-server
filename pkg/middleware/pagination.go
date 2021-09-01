package middleware

import (
	"context"
	"net/http"
	"strconv"

	"quotes/pkg/rest"
	"quotes/pkg/uid"
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

// Pagination middleware attempts to parse pagination query params and passes them in the request context.
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default Pagination data
		pagination := PaginationContext{
			Limit: DefaultLimit,
		}

		// Parse Pagination Params from the request
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

		// 422
		if pagination.Limit < 1 || pagination.Limit > MaxLimit {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		ctx := context.WithValue(r.Context(), PaginationContextKey, &pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPaginationContext returns a pagination request context
func GetPaginationContext(request *http.Request) *PaginationContext {
	return request.Context().Value(PaginationContextKey).(*PaginationContext)
}
