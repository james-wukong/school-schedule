package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiters = make(map[string]*rate.Limiter)

func (m *Manager) RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if _, exists := limiters[ip]; !exists {
			limiters[ip] = rate.NewLimiter(rate.Limit(5), 10)
		}
		limiter := limiters[ip]

		if !limiter.Allow() {
			// utils.ReturnErrorResponse(c, http.StatusTooManyRequests, "Too Many Requests", nil)
			m.log.Warn().Str("ip", ip).Msg("Rate limit exceeded")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
			})
			return
		}
		c.Next()
	}
}
