package handlers

import (
	"TODO/database/dbhelper"
	"TODO/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTodos(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")

	limit, offset := utils.SetPagination(c)

	todos, err := dbhelper.GetAllTodos(status, search, limit, offset)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func SuspendUser(c *gin.Context) {

	currentUser, ok := utils.GetUserFromContext(c)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	targetUser := c.Param("id")
	if targetUser == "" {
		c.JSON(400, gin.H{"error": "user id required"})
		return
	}

	if targetUser == currentUser.ID {
		c.JSON(403, gin.H{"error": "that's suicidal my friend"})
		return
	}

	// fmt.Print(currentUser.IsSuspended)
	// if currentUser.IsSuspended {
	// 	c.JSON(400, gin.H{"error": "user already suspended"})
	// 	return
	// }

	err := dbhelper.SuspendUser(targetUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to suspend user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "user suspended successfully",
	})
}

func UnsuspendUser(c *gin.Context) {

	// currentUser, ok := utils.GetUserFromContext(c)
	// if !ok {
	// 	c.JSON(401, gin.H{"error": "unauthorized"})
	// 	return
	// }

	userID := c.Param("id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user id required"})
		return
	}

	// fmt.Print(currentUser.IsSuspended)
	// if !currentUser.IsSuspended {
	// 	c.JSON(400, gin.H{"error": "user already unsuspended"})
	// 	return
	// }

	err := dbhelper.UnsuspendUser(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to unsuspend user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "user unsuspended successfully",
	})
}

func FetchAllUsers(c *gin.Context) {

	limit, offset := utils.SetPagination(c)

	users, err := dbhelper.FetchAllUsers(limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}
