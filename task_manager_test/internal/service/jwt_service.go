package service

import (
	"errors"
	"task_manager_test/internal/usecase"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtService implements the usecase.JWTService interface.
type jwtService struct{ secret []byte }

// This compile-time check ensures that *jwtService satisfies the usecase.JWTService interface.
var _ usecase.IJWTService = (*jwtService)(nil)

// NewJWTService constructs a new JWTService instance with the provided HMAC secret.
func NewJWTService(secret string) usecase.IJWTService {
	return &jwtService{secret: []byte(secret)}
}

// GenerateToken creates a JWT signed with HS256, containing username, role, and expiration (24h).
func (j *jwtService) GenerateToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken parses and verifies a token string, returning claims if valid.
func (j *jwtService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
