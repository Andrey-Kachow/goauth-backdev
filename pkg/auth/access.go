package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const AccessSecretKey = "accessSecretKey"

const AccessTokenExiparionInMinutes = 15

func GenerateAccessToken(userGUID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": userGUID,
		"exp":  time.Now().Add(time.Minute * AccessTokenExiparionInMinutes).Unix(),
	})
	return token.SignedString([]byte(AccessSecretKey))
}
