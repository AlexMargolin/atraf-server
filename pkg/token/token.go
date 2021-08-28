package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"

	"quotes/pkg/uid"
)

type Claims = jwt.StandardClaims

func New(secret string, subject uid.UID) (string, error) {
	if secret == "" {
		return "", errors.New("empty token secret provided")
	}

	claims := Claims{
		Subject:   uid.ToString(subject),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	}

	// Unsigned Token
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Signed Token
	token, err := unsignedToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func Verify(secret string, unverifiedToken string) (Claims, error) {
	token, err := jwt.ParseWithClaims(unverifiedToken, &Claims{}, signingSecret(secret))

	// Validate token
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
