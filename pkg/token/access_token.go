package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"atraf-server/pkg/uid"
)

const (
	AccessTokenValidFor = time.Minute * 15
)

var AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

type AccessTokenCustomClaims struct {
	Active    bool    `json:"active"`
	AccountId uid.UID `json:"account_id"`
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	AccessTokenCustomClaims
}

func NewAccessToken(claims AccessTokenCustomClaims) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(AccessTokenValidFor).Unix(),
		},
		AccessTokenCustomClaims: claims,
	})

	return unsignedToken.SignedString([]byte(AccessTokenSecret))
}

func VerifyAccessToken(unverifiedToken string) (AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(unverifiedToken, &AccessTokenClaims{}, signingSecret(AccessTokenSecret))
	if err != nil {
		return AccessTokenClaims{}, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return AccessTokenClaims{}, err
	}

	return *claims, nil
}
