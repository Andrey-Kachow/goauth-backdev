package auth

import (
	"testing"
	"time"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/app"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const sampleUserGUID string = "12345"
const sampleClientIP string = "192.168.1.1"
const sampleNewClientIP string = "192.168.1.2"
const sampleUserEmail string = "name1@example.com"

// Mock implementation of NotificationService
type MockNotificationService struct {
	EmailSent     bool
	SentUserEmail string
	SentClientIP  string
}

func (m *MockNotificationService) SendWarning(userEmail string, clientIP string) error {
	m.EmailSent = true
	m.SentUserEmail = userEmail
	m.SentClientIP = clientIP
	return nil
}

var MockAppContext = app.ApplicationContext{
	NotificationService: &MockNotificationService{},
	TokenDB: &MockTokenDB{
		ShouldError: false,
	},
}

// Test function for ValidateAccessTokenClaims
func TestValidateAccessTokenClaims(t *testing.T) {
	accessToken, err := GenerateAccessToken(sampleUserGUID, sampleClientIP, sampleUserEmail)
	assert.NoError(t, err)

	returnedGUID, _, err := ValidateAccessTokenClaims(accessToken, sampleNewClientIP, sampleUserEmail)

	assert.NoError(t, err)
	assert.Equal(t, sampleUserGUID, returnedGUID, "Expected GUID to match")
}

func TestGeneratePair(t *testing.T) {
	mockDB := &MockTokenDB{
		ShouldError: false,
	}
	MockAppContext.TokenDB = mockDB

	accessToken, refreshToken, err := GeneratePair(sampleUserGUID, sampleClientIP, sampleUserEmail, MockAppContext)

	assert.NoError(t, err, "Expected no error from GeneratePair")
	assert.NotEmpty(t, accessToken, "Expected access token to be generated")
	assert.NotEmpty(t, refreshToken, "Expected refresh token to be generated")

	assert.Equal(t, sampleUserGUID, mockDB.SavedUserGUID, "Expected user GUID to match")
	assert.NotEmpty(t, mockDB.SavedHashedTokenHash, "Expected hashed refresh token to be saved in the database")
}

func TestGeneratePair_SaveError(t *testing.T) {
	MockAppContext.TokenDB = &MockTokenDB{
		ShouldError: true,
	}

	accessToken, refreshToken, err := GeneratePair(sampleUserGUID, sampleClientIP, sampleUserEmail, MockAppContext)

	assert.Error(t, err, "Expected an error due to database save failure")
	assert.Empty(t, accessToken, "Expected access token to be empty due to error")
	assert.Empty(t, refreshToken, "Expected refresh token to be empty due to error")
}

func TestGeneratePair_IPChange(t *testing.T) {
	mockDB := &MockTokenDB{
		ShouldError: false,
		SavedIP:     sampleClientIP,
		SavedEmail:  sampleUserEmail,
	}
	mockNotificationService := &MockNotificationService{
		EmailSent:     false,
		SentUserEmail: "",
	}
	MockAppContext.TokenDB = mockDB
	MockAppContext.NotificationService = mockNotificationService

	accessToken, refreshToken, err := GeneratePair(sampleUserGUID, sampleNewClientIP, sampleUserEmail, MockAppContext)

	assert.True(t, mockNotificationService.EmailSent)
	assert.NotEmpty(t, mockNotificationService.SentUserEmail)

	assert.NoError(t, err, "Expected no error from GeneratePair despite the new IP address")
	assert.NotEmpty(t, accessToken, "Expected access token to be generated despite the new IP address")
	assert.NotEmpty(t, refreshToken, "Expected refresh token to be generated despite the new IP address")

	assert.Equal(t, sampleUserGUID, mockDB.SavedUserGUID, "Expected user GUID to match despite the new IP address")
	assert.NotEmpty(t, mockDB.SavedHashedTokenHash, "Expected hashed refresh token to be saved in the database despite the new IP address")
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
