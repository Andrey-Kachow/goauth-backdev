package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

const anotherSampleUserGUID string = "67890"

// TestGenerateRefreshToken tests that GenerateRefreshToken generates a valid refresh token
func TestGenerateRefreshToken(t *testing.T) {

	refreshTokenString, hashedRefreshToken, err := GenerateRefreshToken(anotherSampleUserGUID)

	assert.NoError(t, err, "Expected no error generating the refresh token")
	assert.NotEmpty(t, refreshTokenString, "Expected a non-empty refresh token string")
	assert.NotEmpty(t, hashedRefreshToken, "Expected a non-empty hashed refresh token")

	err = bcrypt.CompareHashAndPassword([]byte(hashedRefreshToken), []byte(refreshTokenString))
	assert.NoError(t, err, "Expected the hashed refresh token to match the original token")
}
