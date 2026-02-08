package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authentication middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header or cookie
		token := c.GetHeader("Authorization")
		if token == "" {
			cookie, err := c.Cookie("auth_token")
			if err == nil {
				token = cookie
			}
		}

		// Remove "Bearer " prefix if present
		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Validate token
		claims, err := validateJWT(token)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("isAdmin", claims.IsAdmin)

		c.Next()
	}
}

// Admin middleware - requires user to be admin
func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Rate limiting middleware
func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement rate limiting
		// Using redis or in-memory cache
		c.Next()
	}
}
