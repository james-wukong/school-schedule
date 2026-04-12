package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func (m *Manager) RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Convert the recovered value 'r' to an error type
				var recoveredErr error
				switch v := err.(type) {
				case error:
					recoveredErr = v // It's already an error
				case string:
					recoveredErr = fmt.Errorf("panic: %s", v) // Convert string to error
				default:
					recoveredErr = fmt.Errorf("panic: %v", v) // Convert any other type to error
				}
				// Log the panic
				m.log.Error().
					Interface("panic_value", recoveredErr). // The value passed to panic()
					Bytes("stack_trace", debug.Stack()).    // Get the full stack trace
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("Panic recovered in handler")

				// Depending on the mode, you might choose to return different responses.
				if gin.Mode() == gin.DebugMode {
					// In debug mode, send a more informative error to the client
					// 1. Write custom JSON response and 2. Abort the chain
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error":   "Internal Server Error - Debug Mode",
						"message": fmt.Sprintf("Panic: %v", recoveredErr),
						"stack":   string(debug.Stack()),
					})
					return // Ensure the rest of THIS function stops executing
				} else {
					// In release mode, send a generic error to the client
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": fmt.Sprintf("Internal Server Error - %s", gin.ReleaseMode),
					})
					return
				}
			}
		}()
		// Proceed with the next handler
		c.Next()
	}
}
