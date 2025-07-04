package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expectedApiKey := os.Getenv("API_KEY")
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" || expectedApiKey == "" || apiKey != expectedApiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key inválida ou ausente"})
			return
		}
		c.Next()
	}
}
