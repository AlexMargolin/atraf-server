package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

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

func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 401
		unverifiedToken, err := BearerToken(r)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		// 401
		claims, err := token.Verify(unverifiedToken)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		// 500
		accountId, err := uid.FromString(claims.Subject)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		sessionContext := &SessionContext{
			AccountId: accountId,
		}

		ctx := context.WithValue(r.Context(), SessionContextKey, sessionContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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

func GetSessionContext(request *http.Request) *SessionContext {
	return request.Context().Value(SessionContextKey).(*SessionContext)
}
