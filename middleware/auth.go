package middleware

import (
	"database/sql"
	"net/http"

	"TODO/database"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionID := c.GetHeader("session_id")
		if sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"sessiod id":"session dalo na",
			})
			return
		}

		var userID string

		query := `
			SELECT user_id	FROM user_session WHERE id = $1 AND archived_at IS NULL`

		err := database.Todo.QueryRow(query, sessionID).Scan(&userID)
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
		c.Set("user_id", userID)
		c.Set("session_id",sessionID)

		c.Next()
	}
}