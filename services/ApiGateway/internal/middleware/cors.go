package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kareemhamed001/e-commerce/pkg/logger"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(allowedOrigins, allowedMethods, allowedHeaders []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed
		allowedOrigin := "*"
		for _, allowed := range allowedOrigins {
			if allowed == "*" || allowed == origin {
				allowedOrigin = allowed
				break
			}
		}

		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", joinStrings(allowedMethods, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Headers", joinStrings(allowedHeaders, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Recovery middleware recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered: %v", err)
				writeJSONError(c, http.StatusInternalServerError, "internal server error")
			}
		}()

		c.Next()
	}
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
