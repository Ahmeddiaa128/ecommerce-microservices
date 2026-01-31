package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeJSONError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error":   http.StatusText(statusCode),
		"message": message,
		"code":    statusCode,
	})
}
