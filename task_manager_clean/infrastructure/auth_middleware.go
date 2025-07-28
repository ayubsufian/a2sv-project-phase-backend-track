package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function that checks for a valid JWT token in the Authorization header.
func AuthMiddleware(jwtSvc JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwtSvc.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}
		// Set username and role
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}
		c.Next()
	}
}

// AdminOnly is a middleware function that restricts access to only users with the "admin" role.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIfc, exists := c.Get("role")
		role, _ := roleIfc.(string)
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
