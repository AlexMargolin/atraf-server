package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
)

const (
	LimitParam  = "limit"
	CursorParam = "cursor"
)

const (
	DefaultLimit = 9
	MaxLimit     = 20
)

type Cursor struct {
	Key   uid.UID   `json:"key"`
	Value time.Time `json:"value"`
}

type PaginationContext struct {
	Limit  int
	Cursor Cursor
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
			if err := DecodeCursor(cursor, &pagination.Cursor); err != nil {
				rest.Error(w, err, http.StatusUnprocessableEntity)
				return
			}
		}

		if pagination.Limit < 1 || pagination.Limit > MaxLimit {
			err := errors.New("invalid pagination params")
			rest.Error(w, err, http.StatusUnprocessableEntity)
			return
		}

		ctx := context.WithValue(r.Context(), PaginationContextKey, &pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DecodeCursor(s string, dest interface{}) error {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &dest); err != nil {
		return err
	}

	return nil
}

func EncodeCursor(p *Cursor) (string, error) {
	str, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(str), nil
}

func GetPaginationContext(request *http.Request) *PaginationContext {
	return request.Context().Value(PaginationContextKey).(*PaginationContext)
}
