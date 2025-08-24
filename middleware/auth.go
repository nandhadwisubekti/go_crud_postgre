package middleware

import (
	"net/http"
	"strings"

	"go-crud-employee/models"
	"go-crud-employee/utils"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager *utils.JWTManager
}

func NewAuthMiddleware(jwtManager *utils.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// RequireAuth middleware validates JWT token from Authorization header
func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"Authorization required",
				"Missing Authorization header",
			))
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"Invalid authorization format",
				"Authorization header must start with 'Bearer '",
			))
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"Token required",
				"Empty token provided",
			))
			c.Abort()
			return
		}

		// Validate token
		claims, err := a.jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"Invalid token",
				err.Error(),
			))
			c.Abort()
			return
		}

		// Store user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		// Continue to next handler
		c.Next()
	}
}

// GetUserFromContext extracts user information from gin context
func GetUserFromContext(c *gin.Context) (*models.JWTClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*models.JWTClaims)
	if !ok {
		return nil, false
	}

	return userClaims, true
}

// GetUserID extracts user ID from gin context
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(int)
	if !ok {
		return 0, false
	}

	return id, true
}