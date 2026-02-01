package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cancellation stops handling if the request context is canceled.
func Cancellation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		select {
		case <-ctx.Done():
			writeJSONError(c, http.StatusServiceUnavailable, "request canceled")
			return
		default:
		}

		c.Next()

		if ctx.Err() != nil && !c.Writer.Written() {
			status := http.StatusServiceUnavailable
			if ctx.Err() == context.DeadlineExceeded {
				status = http.StatusGatewayTimeout
			}
			writeJSONError(c, status, "request canceled")
		}
	}
}
