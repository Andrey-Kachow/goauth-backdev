package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const sampleUserGUID string = "12345"

// TestGenerateAccessToken tests that GenerateAccessToken generates a valid token
func TestGenerateAccessToken(t *testing.T) {
	// Act
	tokenString, err := GenerateAccessToken(sampleUserGUID)

	// Assert
	assert.NoError(t, err, "Expected no error generating the access token")
	assert.NotEmpty(t, tokenString, "Expected a non-empty token string")

	// Parse the token to verify its claims
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(AccessSecretKey), nil
		},
	)

	// Ensure the token is valid
	assert.NoError(t, err, "Expected token to be valid")
	assert.True(t, token.Valid, "Expected the token to be valid")

	// Check the claims in the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		assert.Equal(t, sampleUserGUID, claims["guid"], "Expected GUID to match")

		exp := int64(claims["exp"].(float64)) // JWT encodes numeric claims as float64
		expirationTime := time.Unix(exp, 0)
		expectedExpiration := time.Now().Add(time.Minute * AccessTokenExiparionInMinutes)

		// Allow a small margin for the timing difference in expiration
		assert.WithinDuration(t, expectedExpiration, expirationTime, time.Minute, "Expected expiration to match")
	} else {
		t.Errorf("Expected valid claims in the token")
	}
}
