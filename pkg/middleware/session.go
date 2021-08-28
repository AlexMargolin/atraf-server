package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"quotes/pkg/rest"
	"quotes/pkg/token"
	"quotes/pkg/uid"
)

const (
	AuthTokenHeader = "Authorization"
)

type SessionContext struct {
	AccountId uid.UID
}

type sessionContextKey string

const SessionContextKey sessionContextKey = "SessionCtx"

func Session(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			unverifiedToken, err := BearerToken(r)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized) // 401
				return
			}

			// Verify Token & retrieve claims
			claims, err := token.Verify(secret, unverifiedToken)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized) // 401
				return
			}

			accountId, err := uid.FromString(claims.Subject)
			if err != nil {
				rest.Error(w, http.StatusInternalServerError) // 500
				return
			}

			// Session Context
			sessionContext := &SessionContext{
				AccountId: accountId,
			}

			// Next & Context
			ctx := context.WithValue(r.Context(), SessionContextKey, sessionContext)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// BearerToken retrieves the Auth Token from the Authentication Header
func BearerToken(r *http.Request) (string, error) {
	header := r.Header.Get(AuthTokenHeader)
	if header == "" {
		return "", errors.New("missing auth header")
	}

	splitHeader := strings.Split(header, "Bearer ")
	if len(splitHeader) != 2 {
		return "", errors.New("invalid auth header")
	}

	return splitHeader[1], nil
}

// GetSessionContext returns session request context
func GetSessionContext(request *http.Request) *SessionContext {
	return request.Context().Value(SessionContextKey).(*SessionContext)
}
