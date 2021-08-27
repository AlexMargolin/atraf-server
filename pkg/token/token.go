package token

import (
	"time"

	"github.com/golang-jwt/jwt"

	"quotes/pkg/uid"
)

type Claims = jwt.StandardClaims

func New(secret string, subject uid.UID) (string, error) {
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
	var claims Claims

	// Parse Token
	token, err := jwt.ParseWithClaims(unverifiedToken, claims, signingSecret(secret))

	// Validate token
	claims, ok := token.Claims.(Claims)
	if !ok || !token.Valid {
		return claims, err
	}

	return claims, nil
}

func signingSecret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}
