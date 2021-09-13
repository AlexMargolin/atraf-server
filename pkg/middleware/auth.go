package middleware

import (
	"context"
	"net/http"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/session"
	"atraf-server/pkg/uid"
)

type authContextKey string

const (
	AuthContextKey authContextKey = "AuthCtx"
)

type AuthContext struct {
	AccountId uid.UID
}

func Authenticate(activated bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := session.ReadCookie(r)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			if activated != cookie.AccountActive {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthContextKey, &AuthContext{
				cookie.AccountId,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthContext(request *http.Request) *AuthContext {
	return request.Context().Value(AuthContextKey).(*AuthContext)
}
