package token

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/acanoe/newsbytes-api-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId uint) (string, error) {
	apiSecret := utils.GetEnv("API_SECRET", "secret")                        // API secret
	tokenLifespan, err := strconv.Atoi(utils.GetEnv("TOKEN_LIFESPAN", "30")) // token lifespan in minutes
	if err != nil {
		return "", err
	}

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

func ExtractToken(c *gin.Context) string {
	// if the token is put in the URL in the form of http://url.domain/path?token=...
	token := c.Query("token")
	if token != "" {
		return token
	}

	// if the token is put in the header in the form of Authorization: Bearer token
	authHeader := c.Request.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) == 2 {
		return strings.Split(authHeader, " ")[1]
	}

	return ""
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	apiSecret := utils.GetEnv("API_SECRET", "secret")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(apiSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil

}

func ValidateToken(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := ParseToken(tokenString)
	if err != nil {
		return err
	}

	return nil
}

func GetIDFromToken(c *gin.Context) (uint, error) {
	tokenString := ExtractToken(c)
	token, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, nil
	}

	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(uid), nil
}
