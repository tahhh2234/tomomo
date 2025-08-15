package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tahhh2234/tomomo/backend/internal/auth"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.ParseAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired access token"})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
