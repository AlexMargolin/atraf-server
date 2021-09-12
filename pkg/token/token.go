package token

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

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
