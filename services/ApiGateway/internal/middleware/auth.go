package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	customJWT "github.com/kareemhamed001/e-commerce/pkg/jwt"
	"github.com/kareemhamed001/e-commerce/pkg/logger"
)

type contextKey string

const (
	UserClaimsKey contextKey = "userClaims"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtManager *customJWT.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			writeJSONError(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.Verify(tokenString)
		if err != nil {
			logger.Errorf("JWT validation failed: %v", err)
			writeJSONError(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		// Add claims to context
		ctx := context.WithValue(c.Request.Context(), UserClaimsKey, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware(jwtManager *customJWT.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				claims, err := jwtManager.Verify(tokenString)
				if err == nil {
					ctx := context.WithValue(c.Request.Context(), UserClaimsKey, claims)
					c.Request = c.Request.WithContext(ctx)
				}
			}
		}
		c.Next()
	}
}

// RequireRole checks if user has required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Request.Context().Value(UserClaimsKey).(*customJWT.UserClaims)
		if !ok {
			writeJSONError(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		logger.Infof("User ID %d with role %s is accessing %s", claims.UserID, claims.Role, c.Request.URL.Path)
		hasRole := false
		for _, role := range roles {
			if claims.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			writeJSONError(c, http.StatusForbidden, "insufficient permissions")
			logger.Info("forbidden access attempt by user ID ", claims.UserID)
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserClaims retrieves user claims from context
func GetUserClaims(ctx context.Context) (*customJWT.UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*customJWT.UserClaims)
	return claims, ok
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) (uint, bool) {
	claims, ok := GetUserClaims(ctx)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}

// GetUserRole retrieves user role from context
func GetUserRole(ctx context.Context) (string, bool) {
	claims, ok := GetUserClaims(ctx)
	if !ok {
		return "", false
	}
	return claims.Role, true
}
