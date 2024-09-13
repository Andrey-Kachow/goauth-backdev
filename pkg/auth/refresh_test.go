package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

const anotherSampleUserGUID string = "67890"

type MockTokenDB struct {
	SavedUserGUID        string
	SavedHashedTokenHash string
	ShouldError          bool
}

func (mockDB *MockTokenDB) SaveHashedRefreshToken(userGUID string, refreshTokenHash string) error {
	if mockDB.ShouldError {
		return errors.New("Database write error")
	}
	mockDB.SavedUserGUID = userGUID
	mockDB.SavedHashedTokenHash = refreshTokenHash
	return nil
}

func (mockDB *MockTokenDB) FetchHashedRefreshTokenFromDB(userGUID string) (string, error) {
	if mockDB.ShouldError {
		return "", errors.New("Database read error")
	}
	return mockDB.SavedHashedTokenHash, nil
}

func (mockDB *MockTokenDB) GetEmailAddressFromGUID(userGUID string) (string, error) {
	return "", nil
}

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

func TestValidateRefreshTokenAndPassword(t *testing.T) {
	refreshToken, _, _ := GenerateRefreshToken(anotherSampleUserGUID)
	password := createPasswordFromRefreshToken(refreshToken)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockDB := &MockTokenDB{
		SavedHashedTokenHash: string(hashedPassword),
		ShouldError:          false,
	}

	returnedGUID, err := ValidateRefreshTokenAndPassword(refreshToken, mockDB)
	assert.NoError(t, err)
	assert.Equal(t, anotherSampleUserGUID, returnedGUID, "Expected GUID to match")
}

func TestValidateRefreshTokenAndPassword_DBError(t *testing.T) {
	refreshToken, _, err := GenerateRefreshToken(anotherSampleUserGUID)
	assert.NoError(t, err)

	mockDB := &MockTokenDB{
		ShouldError: true, // Simulate a database error
	}

	returnedGUID, err := ValidateRefreshTokenAndPassword(refreshToken, mockDB)
	assert.Error(t, err, "Expected an error due to database failure")
	assert.Empty(t, returnedGUID, "Expected no GUID to be returned")
}
