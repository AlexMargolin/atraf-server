package middleware

import (
	"context"
	"net/http"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/token"
	"atraf-server/pkg/uid"
)

type sessionContextKey string

const (
	AuthTokenHeader                     = "Authorization"
	SessionContextKey sessionContextKey = "SessionCtx"
)

type SessionContext struct {
	AccountId uid.UID
}

func Session(active bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			unverifiedToken, err := token.FromHeader(r, AuthTokenHeader)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			claims, err := token.VerifyAccessToken(unverifiedToken)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			if active && claims.Active == false {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), SessionContextKey, &SessionContext{
				claims.AccountId,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSessionContext(request *http.Request) *SessionContext {
	return request.Context().Value(SessionContextKey).(*SessionContext)
}
