package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"atraf-server/pkg/uid"
)

const (
	ResetTokenValidFor = time.Minute * 5
)

var ResetTokenSecret = os.Getenv("RESET_TOKEN_SECRET")

type ResetTokensCustomClaims struct {
	AccountId uid.UID `json:"account_id"`
}

type ResetTokenClaims struct {
	jwt.StandardClaims
	ResetTokensCustomClaims
}

func NewResetToken(claims ResetTokensCustomClaims) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, ResetTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(ResetTokenValidFor).Unix(),
		},
		ResetTokensCustomClaims: claims,
	})

	return unsignedToken.SignedString([]byte(ResetTokenSecret))
}

func VerifyResetToken(unverifiedToken string) (ResetTokenClaims, error) {
	token, err := jwt.ParseWithClaims(unverifiedToken, &ResetTokenClaims{}, signingSecret(ResetTokenSecret))
	if err != nil {
		return ResetTokenClaims{}, err
	}

	claims, ok := token.Claims.(*ResetTokenClaims)
	if !ok || !token.Valid {
		return ResetTokenClaims{}, err
	}

	return *claims, nil
}
