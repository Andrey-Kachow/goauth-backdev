package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const sampleUserGUID string = "12345"
const sampleClientIP string = "192.168.1.1"
const sampleNewClientIP string = "192.168.1.2"
const sampleUserEmail string = "name1@example.com"

// Mock implementation of NotificationService
type MockNotificationService struct {
	EmailSent    bool
	SentUserGUID string
	SentClientIP string
}

func (m *MockNotificationService) SendWarning(userGUID string, clientIP string) error {
	m.EmailSent = true
	m.SentUserGUID = userGUID
	m.SentClientIP = clientIP
	return nil
}

func (m *MockNotificationService) GetEmailAddressFromGUID(userGUID string) (string, error) {
	return "", nil
}

// Test function for ValidateAccessTokenClaims
func TestValidateAccessTokenClaims_IPChange(t *testing.T) {
	accessToken, err := GenerateAccessToken(sampleUserGUID, sampleClientIP, sampleUserEmail)
	assert.NoError(t, err)

	mockNotificationService := &MockNotificationService{}
	returnedGUID, _, err := ValidateAccessTokenClaims(accessToken, sampleNewClientIP, sampleUserEmail, mockNotificationService)

	assert.NoError(t, err)
	assert.Equal(t, sampleUserGUID, returnedGUID, "Expected GUID to match")

	assert.True(t, mockNotificationService.EmailSent, "Expected email to be sent due to IP change")
	assert.Equal(t, sampleUserGUID, mockNotificationService.SentUserGUID, "Expected email to be sent to the correct user")
	assert.Equal(t, sampleNewClientIP, mockNotificationService.SentClientIP, "Expected email to contain the new client IP")
}

func TestValidateAccessTokenClaims_NoIPChange(t *testing.T) {
	accessToken, err := GenerateAccessToken(sampleUserGUID, sampleClientIP, sampleUserEmail)
	assert.NoError(t, err)

	mockNotificationService := &MockNotificationService{}
	returnedGUID, _, err := ValidateAccessTokenClaims(accessToken, sampleClientIP, sampleUserEmail, mockNotificationService)

	assert.NoError(t, err)
	assert.Equal(t, sampleUserGUID, returnedGUID, "Expected GUID to match")
	assert.False(t, mockNotificationService.EmailSent, "Expected no email to be sent because the IP hasn't changed")
}

func TestGeneratePair(t *testing.T) {
	mockDB := &MockTokenDB{
		ShouldError: false,
	}
	accessToken, refreshToken, err := GeneratePair(sampleUserGUID, sampleClientIP, sampleUserEmail, mockDB)

	assert.NoError(t, err, "Expected no error from GeneratePair")
	assert.NotEmpty(t, accessToken, "Expected access token to be generated")
	assert.NotEmpty(t, refreshToken, "Expected refresh token to be generated")

	assert.Equal(t, sampleUserGUID, mockDB.SavedUserGUID, "Expected user GUID to match")
	assert.NotEmpty(t, mockDB.SavedHashedTokenHash, "Expected hashed refresh token to be saved in the database")
}

func TestGeneratePair_SaveError(t *testing.T) {
	mockDB := &MockTokenDB{
		ShouldError: true,
	}

	accessToken, refreshToken, err := GeneratePair(sampleUserGUID, sampleClientIP, sampleUserEmail, mockDB)

	assert.Error(t, err, "Expected an error due to database save failure")
	assert.Empty(t, accessToken, "Expected access token to be empty due to error")
	assert.Empty(t, refreshToken, "Expected refresh token to be empty due to error")
}

// TestGenerateAccessToken tests that GenerateAccessToken generates a valid token
func TestGenerateAccessToken(t *testing.T) {
	accessTokenString, err := GenerateAccessToken(sampleUserGUID, sampleClientIP, sampleUserEmail)

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
