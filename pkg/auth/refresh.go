package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

func GenerateRefreshToken(userGUID string, userEmail string) (string, string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"guid":  userGUID,
		"email": userEmail,
		"exp":   time.Now().Add(RefreshKeyExpirationDuration).Unix(),
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

// ValidateRefreshTokenAndPassword returns validated userGUID and email
func ValidateRefreshTokenAndPassword(refreshToken string, tokenDB db.TokenDB) (string, string, error) {
	token, err := jwt.Parse(
		refreshToken,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(RefreshSecretKey), nil
		},
	)
	if err != nil {
		return "", "", err
	}

	claims, isOK := token.Claims.(jwt.MapClaims)
	if !(isOK || token.Valid) {
		return "", "", err
	}

	userGUID, ok := claims["guid"].(string)
	if !ok {
		return "", "", errors.New("user GUID not found in the refresh token")
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		return "", "", errors.New("user email not found in the refresh token")
	}

	hashedTokenFromDB, err := tokenDB.FetchHashedRefreshTokenFromDB(userGUID)
	if err != nil {
		return "", "", err
	}

	passwordFromRefreshToken := createPasswordFromRefreshToken(refreshToken)

	err = bcrypt.CompareHashAndPassword([]byte(hashedTokenFromDB), []byte(passwordFromRefreshToken))
	if err != nil {
		return "", "", err
	}

	return userGUID, userEmail, nil
}
