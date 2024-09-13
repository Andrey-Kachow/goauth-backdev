package auth

import (
	"time"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/db"
	"github.com/golang-jwt/jwt/v5"
)

const AccessSecretKey = "accessSecretKey"

const AccessTokenExiparionInMinutes = 15

func GeneratePair(userGUID string, tokenDB db.TokenDB) (string, string, error) {
	accessToken, err := GenerateAccessToken(userGUID)
	if err != nil {
		return "", "", err
	}

	refreshToken, hashedRefreshToken, err := GenerateRefreshToken(userGUID)
	if err != nil {
		return "", "", err
	}

	err = tokenDB.SaveHashedRefreshToken(userGUID, hashedRefreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func GenerateAccessToken(userGUID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid": userGUID,
		"exp":  time.Now().Add(time.Minute * AccessTokenExiparionInMinutes).Unix(),
	})
	return token.SignedString([]byte(AccessSecretKey))
}
