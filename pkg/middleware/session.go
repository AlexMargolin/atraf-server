package middleware

import (
	"context"
	"net/http"

	"atraf-server/domain/users"
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
	UserId    uid.UID
}

func Session(service *users.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			unverifiedToken, err := token.FromHeader(r, AuthTokenHeader)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			claims, err := token.Verify(unverifiedToken)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			accountId, err := uid.FromString(claims.Subject)
			if err != nil {
				rest.Error(w, http.StatusInternalServerError)
				return
			}

			user, err := service.ByAccount(accountId)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), SessionContextKey, &SessionContext{accountId, user.Id})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSessionContext(request *http.Request) *SessionContext {
	return request.Context().Value(SessionContextKey).(*SessionContext)
}
