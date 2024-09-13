package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const sampleUserGUID string = "12345"

func TestGeneratePair(t *testing.T) {
	// Create a mock TokenDB
	mockDB := &MockTokenDB{
		ShouldError: false, // No error during save
	}
	accessToken, refreshToken, err := GeneratePair(sampleUserGUID, mockDB)

	assert.NoError(t, err, "Expected no error from GeneratePair")
	assert.NotEmpty(t, accessToken, "Expected access token to be generated")
	assert.NotEmpty(t, refreshToken, "Expected refresh token to be generated")

	assert.Equal(t, sampleUserGUID, mockDB.SavedUserGUID, "Expected user GUID to match")
	assert.NotEmpty(t, mockDB.SavedHashedTokenHash, "Expected hashed refresh token to be saved in the database")
}

func TestGeneratePair_SaveError(t *testing.T) {

	mockDB := &MockTokenDB{
		ShouldError: true, // Simulate an error during save
	}

	accessToken, refreshToken, err := GeneratePair(sampleUserGUID, mockDB)

	assert.Error(t, err, "Expected an error due to database save failure")
	assert.Empty(t, accessToken, "Expected access token to be empty due to error")
	assert.Empty(t, refreshToken, "Expected refresh token to be empty due to error")
}

// TestGenerateAccessToken tests that GenerateAccessToken generates a valid token
func TestGenerateAccessToken(t *testing.T) {
	accessTokenString, err := GenerateAccessToken(sampleUserGUID)

	assert.NoError(t, err, "Expected no error generating the access token")
	assert.NotEmpty(t, accessTokenString, "Expected a non-empty token string")

	// Parse the token to verify its claims
	//
	token, err := jwt.Parse(
		accessTokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(AccessSecretKey), nil
		},
	)

	assert.NoError(t, err, "Expected token to be valid")
	assert.True(t, token.Valid, "Expected the token to be valid")

	// Check the claims in the token
	//
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		t.Errorf("Expected valid claims in the access token")
		return
	}

	assert.Equal(t, sampleUserGUID, claims["guid"], "Expected GUID to match")

	exp := int64(claims["exp"].(float64))
	expirationTime := time.Unix(exp, 0)
	expectedExpiration := time.Now().Add(time.Minute * AccessTokenExiparionInMinutes)

	// Allow a small margin for the timing difference in expiration
	//
	assert.WithinDuration(t, expectedExpiration, expirationTime, time.Minute, "Expected expiration to match")
}
