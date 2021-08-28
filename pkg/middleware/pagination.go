package middleware

import (
	"context"
	"net/http"
	"strconv"
)

const (
	PageNumParam = "page"
	PerPageParam = "per_page"
)

const (
	DefaultPageNum = 1
	DefaultPerPage = 9
)

type PaginationContext struct {
	PageNum int
	PerPage int
}

type paginationContextKey string

const PaginationContextKey paginationContextKey = "PaginationCtx"

// Pagination middleware attempts to parse pagination query params and passes them in the request context.
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default Pagination data
		pagination := PaginationContext{
			PageNum: DefaultPageNum,
			PerPage: DefaultPerPage,
		}

		// Parse Pagination Params from the request
		pageNum := r.URL.Query().Get(PageNumParam) // Page Number
		perPage := r.URL.Query().Get(PerPageParam) // Page Limit

		if pageNum != "" {
			if pageNum, err := strconv.Atoi(pageNum); err == nil {
				pagination.PageNum = pageNum
			}
		}

		if perPage != "" {
			if perPage, err := strconv.Atoi(perPage); err == nil {
				pagination.PerPage = perPage
			}
		}

		// Next & Context
		ctx := context.WithValue(r.Context(), PaginationContextKey, &pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HasCurrentPage determines whether the current pagination settings would result in anymore data
func HasCurrentPage(pagination *PaginationContext, total int) bool {
	// we subtract 1 in order to be able to always show results
	// at the first page regardless of what perPage value is.
	// since 0*x will be smaller than total. unless total is also 0.
	return (pagination.PageNum-1)*pagination.PerPage < total
}

// HasNextPage determines whether the next pagination settings would result in anymore data
func HasNextPage(pagination *PaginationContext, total int) bool {
	return pagination.PageNum*pagination.PerPage < total
}

// GetPaginationContext returns a pagination request context
func GetPaginationContext(request *http.Request) *PaginationContext {
	return request.Context().Value(PaginationContextKey).(*PaginationContext)
}
