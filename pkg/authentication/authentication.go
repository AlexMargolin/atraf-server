package authentication

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
)

type contextKey string

const (
	AccessTokenCookie = "atcId"
	AccessTokenExpiry = time.Minute * 100 // TODO set at 10 minutes
)

const ContextKey contextKey = "AuthCtx"

type CustomClaims struct {
	AccountId     uid.UID `json:"account_id"`
	AccountActive bool    `json:"account_active"`
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	CustomClaims
}

var AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

func SetCookie(w http.ResponseWriter, accountId uid.UID, accountActive bool) error {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(AccessTokenExpiry).Unix(),
		},
		CustomClaims: CustomClaims{
			AccountId:     accountId,
			AccountActive: accountActive,
		},
	})

	token, err := unsignedToken.SignedString([]byte(AccessTokenSecret))
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    token,
		Path:     "/",
		Secure:   false, // TODO enable
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, cookie)

	return nil
}

func ReadCookie(r *http.Request) (*AccessTokenClaims, error) {
	cookie, err := r.Cookie(AccessTokenCookie)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &AccessTokenClaims{}, signingSecret(AccessTokenSecret))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func signingSecret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}

func Middleware(activated bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := ReadCookie(r)
			if err != nil {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			if activated != claims.AccountActive {
				rest.Error(w, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKey, &claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Context(request *http.Request) *AccessTokenClaims {
	return request.Context().Value(ContextKey).(*AccessTokenClaims)
}
