package middleware

import (
	"TODO/database/dbhelper"
	"TODO/models"
	"TODO/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Extract session_id from token
		sessionID := claims.SessionID
		if sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid session in token",
			})
			return
		}

		// Validate session from DB
		userID, err := dbhelper.GetUserIDBySession(sessionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired session",
			})
			return
		}

		// Cross-check token vs DB
		if userID != claims.UserID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token mismatch",
			})
			return
		}

		user, err := dbhelper.GetUserByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "user not found",
			})
			return
		}

		// Suspension check
		if user.IsSuspended {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "account suspended",
			})
			return
		}

		// var userID string

		// userID, err := database.GetUserIDBySession(sessionID)
		// if err != nil {

		// 	if err == sql.ErrNoRows {
		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 			"error": "invalid or expired session",
		// 		})
		// 		return
		// 	}

		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		// 		"error": "database error",
		// 	})
		// 	return
		// }

		// var auth model.AuthContext
		// auth.UserID = userID
		// auth.SessionID = sessionID
		// c.Set("auth", auth)
		// c.Set("user_id", userID)
		c.Set("session_id", sessionID)
		c.Set("user", user)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		u, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		user := u.(models.User)

		if user.Role != "admin" {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
