package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryWithLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %+v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
				c.Abort()
			}
		}()
		c.Next()
	}
}
