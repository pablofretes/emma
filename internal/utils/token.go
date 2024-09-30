package utils

import (
	"fmt"
	"time"

	"emma/configs"

	"github.com/golang-jwt/jwt"
)

var hmacSecret = []byte(configs.GetConfig().JWT_SECRET)
var expiryTime = configs.GetConfig().JWT_EXPIRY_TIME

type SignTokenClaims struct {
	Id       string
	Username string
	Role     string
}

func Authenticate(token string) (result *jwt.Token, err error) {
	result, err = jwt.Parse(token, func(result *jwt.Token) (interface{}, error) {
		if _, ok := result.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", result.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return result, err
}

func SignToken(claims SignTokenClaims) (tokenString string, err error) {
	duration, err := time.ParseDuration(expiryTime)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       claims.Id,
		"username": claims.Username,
		"role":     claims.Role,
		"StandardClaims": jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	})
	tokenString, err = token.SignedString(hmacSecret)
	return tokenString, err
}
