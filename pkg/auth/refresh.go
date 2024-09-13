package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const RefreshSecretKey = "refreshSecretKey"
const RefreshKeyExpirationDuration time.Duration = time.Hour * 24 * 7

// createPasswordFromRefreshToken creates a shorter password for using sha256 algorithm
func createPasswordFromRefreshToken(inputString string) string {
	sha256Hasher := sha256.New()
	sha256Hasher.Write([]byte(inputString))
	hashedRefreshTokenBytes := sha256Hasher.Sum(nil) // This is 32 bytes long
	return hex.EncodeToString(hashedRefreshTokenBytes)
}

func GenerateRefreshToken(userGUID string) (string, string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"guid": userGUID,
		"exp":  time.Now().Add(RefreshKeyExpirationDuration).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(RefreshSecretKey))
	if err != nil {
		return "", "", err
	}

	hashedRefreshTokenPasswordGenerationKey := createPasswordFromRefreshToken(refreshTokenString)

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(hashedRefreshTokenPasswordGenerationKey), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshTokenString, string(hashedRefreshToken), nil
}

func ValidateRefreshTokenAndPassword(refreshToken string, tokenDB *db.TokenDB) (string, error) {
	token, err := jwt.Parse(
		refreshToken,
		func(token *jwt.Token) (interface{}, error) {
			return RefreshSecretKey, nil
		},
	)
	if err != nil {
		return "", err
	}

	claims, isOK := token.Claims.(jwt.MapClaims)
	if !(isOK || token.Valid) {
		return "", err
	}

	userGUID := claims["guid"].(string)
	hashedTokenFromDB, err := tokenDB.FetchHashedRefreshTokenFromDB(userGUID)
	if err != nil {
		return "", err
	}

	passwordFromRefreshToken := createPasswordFromRefreshToken(refreshToken)

	err = bcrypt.CompareHashAndPassword([]byte(hashedTokenFromDB), []byte(passwordFromRefreshToken))
	if err != nil {
		return "", err
	}

	return userGUID, nil
}
