package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId uint) (string, error) {
	tokenLifespan := 30   // token lifespan in minutes
	apiSecret := "secret" // API secret

	expirationTime := time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()

	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userId,
		"exp":        expirationTime,
	}

	signingMethod := jwt.SigningMethodHS256
	token := jwt.NewWithClaims(signingMethod, claims)

	return token.SignedString([]byte(apiSecret))
}
