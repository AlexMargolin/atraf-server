package session

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"atraf-server/pkg/uid"
)

const (
	ATCookieName = "token"
	ATValidFor   = time.Minute * 10
)

var secret = os.Getenv("ACCESS_TOKEN_SECRET")

type Data struct {
	AccountId     uid.UID `json:"account_id"`
	AccountActive bool    `json:"account_active"`
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	Data
}

func SetCookie(w http.ResponseWriter, accountId uid.UID, accountActive bool) error {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(ATValidFor).Unix(),
		},
		Data: Data{
			AccountId:     accountId,
			AccountActive: accountActive,
		},
	})

	token, err := unsignedToken.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     ATCookieName,
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
	cookie, err := r.Cookie(ATCookieName)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &AccessTokenClaims{}, signingSecret(secret))
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
