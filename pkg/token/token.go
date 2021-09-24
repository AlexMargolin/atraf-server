package token

import (
	"github.com/golang-jwt/jwt/v4"
)

func signingSecret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}
