package middleware

import (
	database "TODO/database/dbhelper"
	model "TODO/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionID := c.GetHeader("session_id")
		if sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"sessiod id": "session dalo na",
			})
			return
		}

		var userID string

		userID, err := database.GetUserIDBySession(sessionID)
		if err != nil {

			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "invalid or expired session",
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "database error",
			})
			return
		}

		var auth model.AuthContext
		auth.UserID = userID
		auth.SessionID = sessionID
		c.Set("auth", auth)

		// c.Set("user_id", userID)
		// c.Set("session_id", sessionID) //set this up in object

		c.Next()
	}
}
