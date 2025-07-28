package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService defines methods for generating and validating JWT tokens.
type JWTService interface {
	GenerateToken(username, role string) (string, error)
	ValidateToken(tokenStr string) (jwt.MapClaims, error)
}

// jwtService implements JWTService using a secret key for HMAC signing.
type jwtService struct{ secret []byte }

// NewJWTService constructs a new JWTService instance with the provided HMAC secret.
func NewJWTService(secret []byte) JWTService {
	return &jwtService{secret}
}

// GenerateToken creates a JWT signed with HS256, containing username, role, and expiration (24h).
func (s *jwtService) GenerateToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken parses and verifies a token string, returning claims if valid.
func (s *jwtService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
