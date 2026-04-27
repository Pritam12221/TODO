package handlers

import (
	"TODO/database/dbhelper"
	// "TODO/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllTodos(c *gin.Context) {

	// auth, ok := utils.GetAuth(c)
	// if !ok {
	// 	c.JSON(401, gin.H{"error": "unauthorized"})
	// 	return
	// }

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

	// auth, ok := utils.GetAuth(c)
	// if !ok {
	// 	c.JSON(401, gin.H{"error": "unauthorized"})
	// 	return
	// }

	// if auth.Role != "admin" {
	// 	c.JSON(403, gin.H{"error": "forbidden"})
	// 	return
	// }

	err := dbhelper.SuspendUser(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to suspend user"})
		return
	}

	c.JSON(200, gin.H{"message": "user suspended successfully"})
}
