package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

const AccessSecretKey = "accessSecretKey"

const AccessTokenExiparionInMinutes = 15

// GeneratePair generates and returns access token and refresh token plus any error
func GeneratePair(userGUID string, clientIP string, userEmail string, appContext app.ApplicationContext) (string, string, error) {
	userData, err := appContext.TokenDB.FetchUserData(userGUID)
	if err == nil {
		fmt.Printf("Comparing client IPs: %s and %s", userData.RecentIP, clientIP)
		if userData.RecentIP != clientIP {
			appContext.NotificationService.SendWarning(userEmail, clientIP)
		}
	}

	accessToken, err := GenerateAccessToken(userGUID, clientIP, userEmail)
	if err != nil {
		return "", "", err
	}

	refreshToken, hashedRefreshToken, err := GenerateRefreshToken(userGUID, userEmail)
	if err != nil {
		return "", "", err
	}

	err = appContext.TokenDB.SaveUserData(userGUID, userEmail, clientIP, hashedRefreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// GenerateAccessToken generates new access
func GenerateAccessToken(userGUID string, clientIP string, userEmail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"guid":  userGUID,
		"ip":    clientIP,
		"email": userEmail,
		"exp":   time.Now().Add(time.Minute * AccessTokenExiparionInMinutes).Unix(),
	})
	return token.SignedString([]byte(AccessSecretKey))
}

// ValidateAccessTokenClaims checks the guid and cliend IP and sends notification email in case when IP has changed
func ValidateAccessTokenClaims(accessToken string, currentClientIP string, providedUserEmail string) (string, string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(AccessSecretKey), nil
	})
	if err != nil {
		return "", "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return "", "", errors.New("invalid token claims")
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		return "", "", errors.New("user email not gound in token")
	}
	if userEmail != providedUserEmail {
		return "", "", errors.New("provided email does not match the email found in token")
	}

	userGUID, ok := claims["guid"].(string)
	if !ok {
		return "", "", errors.New("user GUID not found in token")
	}
	return userGUID, userEmail, nil
}
