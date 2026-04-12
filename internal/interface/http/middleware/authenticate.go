package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Manager) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			m.log.Warn().Msgf("Missing token in request path: %s", c.Request.URL.Path)

			// 1. Write custom JSON response and 2. Abort the chain
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authentication token is required",
			})
			return // Ensure the rest of THIS function stops executing
		}

		// If valid, just call Next()
		// TODO: Validate the token and set user info in context, against cache
		c.Next()
	}
}
