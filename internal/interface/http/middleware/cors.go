package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Manager) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE, UPDATE",
		)
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Origin, Content-Type, api_key, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		)
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Cache-Control", "no-cache")

		if c.Request.Method == "OPTIONS" {
			// log.Info().Msgf("OPTIONS method is allowed: %v \n", c.Request.Method)
			m.log.Info().Str("ip", c.ClientIP()).Msg("OPTIONS method is allowed")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
