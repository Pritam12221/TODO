package handlers

import (
	"TODO/database/dbhelper"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllTodos(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	todos, err := dbhelper.GetAllTodos(status, search, limit, offset)
	if err != nil {

		fmt.Println("er")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func SuspendUser(c *gin.Context) {

	userID := c.Param("id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user id required"})
		return
	}

	err := dbhelper.SuspendUser(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to suspend user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "user suspended successfully",
	})
}

func UnsuspendUser(c *gin.Context) {

	userID := c.Param("id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user id required"})
		return
	}

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

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	users, err := dbhelper.FetchAllUsers(limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}
