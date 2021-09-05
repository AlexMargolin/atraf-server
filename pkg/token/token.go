package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"atraf-server/pkg/uid"
)

var AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

type Claims = jwt.StandardClaims

func New(subject uid.UID) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		Subject:   uid.ToString(subject),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})

	return unsignedToken.SignedString([]byte(AccessTokenSecret))
}

func Verify(unverifiedToken string) (Claims, error) {
	token, err := jwt.ParseWithClaims(unverifiedToken, &Claims{}, signingSecret(AccessTokenSecret))

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, err
	}

	return *claims, nil
}

func signingSecret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}
