package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const anotherSampleUserGUID string = "67890"

// TestGenerateRefreshToken tests that GenerateRefreshToken generates a valid refresh token
func TestGenerateRefreshToken(t *testing.T) {

	refreshTokenString, hashedRefreshToken, err := GenerateRefreshToken(anotherSampleUserGUID)

	assert.NoError(t, err, "Expected no error generating the refresh token")
	assert.NotEmpty(t, refreshTokenString, "Expected a non-empty refresh token string")
	assert.NotEmpty(t, hashedRefreshToken, "Expected a non-empty hashed refresh token")

	// Parse the refresh token to verify its claims
	//
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(RefreshSecretKey), nil
	})

	assert.NoError(t, err, "Expected refresh token to be valid")
	assert.True(t, token.Valid, "Expected the refresh token to be valid")

	// Check the claims in the token
	//
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		t.Errorf("Expected valid claims in the refresh token")
		return
	}

	assert.Equal(t, anotherSampleUserGUID, claims["guid"], "Expected GUID to match")

	exp := int64(claims["exp"].(float64))
	expirationTime := time.Unix(exp, 0)
	expectedExpiration := time.Now().Add(RefreshKeyExpirationDuration)

	assert.WithinDuration(t, expectedExpiration, expirationTime, time.Minute, "Expected expiration to match")
}
