package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const RefreshSecretKey = "refreshSecretKey"
const RefreshKeyExpirationDuration time.Duration = time.Hour * 24 * 7

func GenerateRefreshToken(userGUID string) (string, string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"guid": userGUID,
		"exp":  time.Now().Add(RefreshKeyExpirationDuration).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(RefreshSecretKey))
	if err != nil {
		return "", "", err
	}

	fmt.Println("length is " + fmt.Sprint(len(refreshTokenString)))

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshTokenString), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshTokenString, string(hashedRefreshToken), nil
}

func ValidateRefreshToken(refreshToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		refreshToken,
		func(token *jwt.Token) (interface{}, error) {
			return RefreshSecretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, isOK := token.Claims.(jwt.MapClaims)
	if !(isOK || token.Valid) {
		return nil, err
	}
	return claims, nil
}
