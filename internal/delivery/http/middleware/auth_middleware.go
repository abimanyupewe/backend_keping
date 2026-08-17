package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"backend_keping/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware verifies Supabase JWT Access Tokens
func AuthMiddleware(supabaseJWTSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid Authorization header format. Format must be 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(supabaseJWTSecret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			sub, exists := claims["sub"].(string)
			if !exists {
				response.Error(c, http.StatusUnauthorized, "Invalid token claims: missing sub")
				c.Abort()
				return
			}

			userUUID, err := uuid.Parse(sub)
			if err != nil {
				response.Error(c, http.StatusUnauthorized, "Invalid user ID format in token")
				c.Abort()
				return
			}

			c.Set("user_id", userUUID)
			c.Next()
		} else {
			response.Error(c, http.StatusUnauthorized, "Failed to parse token claims")
			c.Abort()
			return
		}
	}
}

// IngestionSecretMiddleware verifies secret key for background notification ingestion
func IngestionSecretMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientSecret := c.GetHeader("X-Ingestion-Secret")
		if clientSecret == "" || clientSecret != secretKey {
			response.Error(c, http.StatusForbidden, "Invalid or missing X-Ingestion-Secret header")
			c.Abort()
			return
		}
		c.Next()
	}
}
