package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const ApiKeyEsperada = "minha-chave-secreta" // Troque por sua chave fixa

func ApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != ApiKeyEsperada {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key inválida ou ausente"})
			return
		}
		c.Next()
	}
}
