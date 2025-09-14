package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/repositories"
	"golang.org/x/net/context"
)

func JWTAuthMiddleware(userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		blacklisted, err := userRepo.IsTokenBlacklisted(context.Background(), tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check token"})
			c.Abort()
			return
		}
		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is blacklisted"})
			c.Abort()
			return
		}

		c.Next()
	}
}
