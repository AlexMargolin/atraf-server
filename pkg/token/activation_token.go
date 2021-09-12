package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"atraf-server/pkg/uid"
)

const (
	ActivationTokenValidFor = time.Minute * 10
)

var ActivationTokenSecret = os.Getenv("ACTIVATION_TOKEN_SECRET")

type ActivationTokensCustomClaims struct {
	AccountId uid.UID `json:"account_id"`
}

type ActivationTokenClaims struct {
	jwt.StandardClaims
	ActivationTokensCustomClaims
}

func NewActivationToken(claims ActivationTokensCustomClaims) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, ActivationTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(ActivationTokenValidFor).Unix(),
		},
		ActivationTokensCustomClaims: claims,
	})

	return unsignedToken.SignedString([]byte(ActivationTokenSecret))
}

func VerifyActivationToken(unverifiedToken string) (ActivationTokenClaims, error) {
	token, err := jwt.ParseWithClaims(unverifiedToken, &ActivationTokenClaims{}, signingSecret(ActivationTokenSecret))
	if err != nil {
		return ActivationTokenClaims{}, err
	}

	claims, ok := token.Claims.(*ActivationTokenClaims)
	if !ok || !token.Valid {
		return ActivationTokenClaims{}, err
	}

	return *claims, nil
}
