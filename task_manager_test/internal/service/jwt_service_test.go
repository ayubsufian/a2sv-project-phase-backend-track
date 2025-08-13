package service

import (
	"crypto/rand"
	"crypto/rsa"
	"task_manager_test/internal/usecase"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// JWTServiceTestSuite defines the test suite for the JWT service.
type JWTServiceTestSuite struct {
	suite.Suite
	jwtService usecase.IJWTService
	secretKey  string
}

// SetupTest runs before each test in the suite.
func (s *JWTServiceTestSuite) SetupTest() {
	s.secretKey = "a-very-secure-secret-key-for-testing"
	s.jwtService = NewJWTService(s.secretKey)
}

// TestJWTServiceTestSuite is the entry point for the Go test runner.
func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}

// --- Test Cases ---

// TestGenerateAndValidateToken_RoundTripSuccess tests the primary "happy path" scenario.
func (s *JWTServiceTestSuite) TestGenerateAndValidateToken_RoundTripSuccess() {
	// ARRANGE
	username := "testuser"
	role := "admin"

	// ACT - Generate the token
	tokenString, err := s.jwtService.GenerateToken(username, role)

	// ASSERT - Generation
	assert.NoError(s.T(), err, "Token generation should not produce an error")
	assert.NotEmpty(s.T(), tokenString, "Generated token string should not be empty")

	// ACT - Validate the token
	claims, err := s.jwtService.ValidateToken(tokenString)

	// ASSERT - Validation
	assert.NoError(s.T(), err, "Token validation should not produce an error for a valid token")
	assert.NotNil(s.T(), claims, "Claims should not be nil for a valid token")
	assert.Equal(s.T(), username, claims["username"], "Username in claims should match the original")
	assert.Equal(s.T(), role, claims["role"], "Role in claims should match the original")

	// Verify the expiration claim ('exp') is set correctly in the future
	expClaim, ok := claims["exp"].(float64)
	assert.True(s.T(), ok, "Expiration claim should exist and be a number")
	expectedExp := time.Now().Add(24 * time.Hour).Unix()
	assert.InDelta(s.T(), expectedExp, int64(expClaim), 2, "Expiration time should be approximately 24 hours from now")
}

// TestValidateToken_Fails_When_Expired tests the edge case where a token is expired.
func (s *JWTServiceTestSuite) TestValidateToken_Fails_When_Expired() {
	// ARRANGE
	claims := jwt.MapClaims{
		"username": "expireduser",
		"role":     "user",
		"exp":      time.Now().Add(-1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, err := token.SignedString([]byte(s.secretKey))
	assert.NoError(s.T(), err, "Setup: Failed to sign expired token")

	// ACT
	_, err = s.jwtService.ValidateToken(expiredTokenString)

	// ASSERT
	assert.Error(s.T(), err, "Validation should fail for an expired token")
	assert.ErrorContains(s.T(), err, "token has invalid claims: token is expired", "Error message should indicate token expiration")
}

// TestValidateToken_Fails_When_InvalidSignature tests the critical security case where a token was signed with a different secret key.
func (s *JWTServiceTestSuite) TestValidateToken_Fails_When_InvalidSignature() {
	// ARRANGE
	tokenString, err := s.jwtService.GenerateToken("legituser", "user")
	assert.NoError(s.T(), err, "Setup: Failed to generate token")
	invalidService := NewJWTService("this-is-the-wrong-secret")

	// ACT
	_, err = invalidService.ValidateToken(tokenString)

	// ASSERT
	assert.Error(s.T(), err, "Validation should fail for a token with an invalid signature")
	assert.ErrorContains(s.T(), err, "signature is invalid", "Error message should indicate invalid signature")
}

// TestValidateToken_Fails_When_InvalidSigningMethod tests that the service correctly rejects tokens that use an unexpected signing algorithm.
func (s *JWTServiceTestSuite) TestValidateToken_Fails_When_InvalidSigningMethod() {
	// ARRANGE
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(s.T(), err, "Setup: Failed to generate RSA private key")

	claims := jwt.MapClaims{"username": "hacker", "role": "user"}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	maliciousTokenString, err := token.SignedString(privateKey)
	assert.NoError(s.T(), err, "Setup: Failed to sign token with RSA key")

	// ACT
	_, err = s.jwtService.ValidateToken(maliciousTokenString)

	// ASSERT
	assert.Error(s.T(), err, "Validation should fail for a token with an unexpected signing method")
	// Use ErrorContains to check for our specific error, ignoring the library's wrapper text.
	assert.ErrorContains(s.T(), err, "unexpected signing method", "Error should indicate the signing method mismatch")
}

// TestValidateToken_Fails_When_MalformedToken tests how the service handles input that is not a valid JWT formatted string.
func (s *JWTServiceTestSuite) TestValidateToken_Fails_When_MalformedToken() {
	// ARRANGE
	malformedToken := "this.is.not.a.valid.jwt"

	// ACT
	_, err := s.jwtService.ValidateToken(malformedToken)

	// ASSERT
	assert.Error(s.T(), err, "Validation should fail for a malformed token string")
	assert.ErrorContains(s.T(), err, "token is malformed", "Error should indicate a malformed token")
}
