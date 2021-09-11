package token

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Claims = jwt.StandardClaims

func New(secret string, c Claims) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	return unsignedToken.SignedString([]byte(secret))
}

func Verify(secret string, unverifiedToken string) (Claims, error) {
	token, err := jwt.ParseWithClaims(unverifiedToken, &Claims{}, signingSecret(secret))
	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, err
	}

	return *claims, nil
}

func FromHeader(r *http.Request, key string) (string, error) {
	header := r.Header.Get(key)
	if header == "" {
		return "", errors.New("missing auth header")
	}

	splitHeader := strings.Split(header, "Bearer ")
	if len(splitHeader) != 2 {
		return "", errors.New("invalid auth header")
	}

	return splitHeader[1], nil
}

func signingSecret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}
