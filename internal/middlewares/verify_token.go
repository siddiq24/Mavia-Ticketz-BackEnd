package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/siddiq24/Tickitz-DB/internal/utils"
)

func VerifyToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
		c.Abort()
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		log.Println(err)
		return
	}

	// Set user information in context
	c.Set("userID", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("role", claims.Role)

	c.Next()
}
