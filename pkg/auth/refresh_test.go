package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const anotherSampleUserGUID string = "67890"
const anotherSampleClientIP string = "192.168.1.1"
const anotherSampleUserEmail string = "sample@example.com"

type MockTokenDB struct {
	SavedUserGUID        string
	SavedHashedTokenHash string
	SavedEmail           string
	ShouldError          bool
}

func (mockDB *MockTokenDB) SaveUserData(userGUID string, userEmail string, clientIP string, refreshTokenHash string) error {
	if mockDB.ShouldError {
		return errors.New("Database write error")
	}
	mockDB.SavedUserGUID = userGUID
	mockDB.SavedHashedTokenHash = refreshTokenHash
	mockDB.SavedEmail = userEmail
	return nil
}

func (mockDB *MockTokenDB) FetchUserData(userGUID string) (db.UserData, error) {
	if mockDB.ShouldError {
		return db.UserData{}, errors.New("Database read error")
	}
	userDataRet := db.UserData{
		RefreshTokenHash: mockDB.SavedHashedTokenHash,
	}
	return userDataRet, nil
}

func (mockDB *MockTokenDB) GetEmailAddressFromGUID(userGUID string) (string, error) {
	return "", nil
}

// TestGenerateRefreshToken tests that GenerateRefreshToken generates a valid refresh token
func TestGenerateRefreshToken(t *testing.T) {

	refreshTokenString, hashedRefreshToken, err := GenerateRefreshToken(anotherSampleUserGUID, anotherSampleUserEmail)

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

func TestValidateRefreshTokenAndPassword(t *testing.T) {

	MockAppContext.TokenDB = &MockTokenDB{
		ShouldError: false, // Simulate a database error
	}

	_, refreshToken, err := GeneratePair(anotherSampleUserGUID, anotherSampleClientIP, anotherSampleUserEmail, MockAppContext)
	assert.NoError(t, err)

	returnedGUID, returnedEmail, err := ValidateRefreshTokenAndPassword(refreshToken, MockAppContext)
	assert.NoError(t, err)
	assert.Equal(t, anotherSampleUserGUID, returnedGUID, "Expected GUID to match")
	assert.Equal(t, anotherSampleUserEmail, returnedEmail, "Expected email to match")
}

func TestValidateRefreshTokenAndPassword_DBError(t *testing.T) {
	refreshToken, _, err := GenerateRefreshToken(anotherSampleUserGUID, anotherSampleUserEmail)
	assert.NoError(t, err)

	MockAppContext.TokenDB = &MockTokenDB{
		ShouldError: true, // Simulate a database error
	}

	returnedGUID, returnedEmail, err := ValidateRefreshTokenAndPassword(refreshToken, MockAppContext)
	assert.Error(t, err, "Expected an error due to database failure")
	assert.Empty(t, returnedGUID, "Expected no GUID to be returned")
	assert.Empty(t, returnedEmail, "Expected no email to be returned")
}
