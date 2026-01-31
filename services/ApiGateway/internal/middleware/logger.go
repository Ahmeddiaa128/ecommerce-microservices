package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kareemhamed001/e-commerce/pkg/logger"
)

// Logger middleware logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Get request ID from context
		requestID, ok := c.Get("requestID")
		if !ok {
			requestID = "unknown"
		}

		// Log request details
		duration := time.Since(start)
		logger.Infof(
			"[%s] %s %s - Status: %d - Duration: %v - Size: %d bytes",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			c.Writer.Size(),
		)
	}
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add to response header
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Add to context
		ctx := context.WithValue(c.Request.Context(), "requestID", requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("requestID", requestID)

		c.Next()
	}
}
