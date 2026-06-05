package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware provides JWT and API key authentication.
type AuthMiddleware struct {
	jwtSecret []byte
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: []byte(jwtSecret)}
}

// JWTAuth validates JWT tokens.
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return m.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
			c.Abort()
			return
		}

		orgID, err := uuid.Parse(claims["org"].(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid org id in token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("org_id", orgID)
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

// APIKeyAuth validates API keys.
func (m *AuthMiddleware) APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			c.Abort()
			return
		}

		// API key validation is handled by the service layer
		// This middleware just extracts the key for validation
		c.Set("api_key", apiKey)
		c.Next()
	}
}
