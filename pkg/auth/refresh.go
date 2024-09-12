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

	refreshTokenString, err := refreshToken.SignedString(RefreshSecretKey)
	if err != nil {
		return "", "", err
	}

	// Hash the refresh token before storing in DB
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshTokenString), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshTokenString, string(hashedRefreshToken), nil
}

func ValidateToken(token string) (map[string]interface{}, error) {
	// Validate the token signature and expiration
	return nil, nil
}

// Refresh Access Token
func refreshAccessToken(refreshToken string, hashedTokenFromDB string) (string, error) {
	// Validate the refresh token
	err := bcrypt.CompareHashAndPassword([]byte(hashedTokenFromDB), []byte(refreshToken))
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Parse the refresh token to extract claims (like the user's GUID)
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return RefreshSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userGUID := claims["guid"].(string)

		// Generate a new access token
		return GenerateAccessToken(userGUID)
	}

	return "", fmt.Errorf("invalid refresh token")
}
