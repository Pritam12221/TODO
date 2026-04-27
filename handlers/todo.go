package handlers

import (
	"TODO/database/dbhelper"
	model "TODO/models"
	util "TODO/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {

	auth, ok := util.GetAuth(c)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userID := auth.UserID

	var todoReq model.Todo

	if err := c.ShouldBindJSON(&todoReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := dbhelper.CreateTodo(userID, todoReq.Name, todoReq.Description, todoReq.ExpiringAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c *gin.Context) {
	auth, ok := util.GetAuth(c)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userID := auth.UserID

	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "todo id is required"})
		return
	}

	var req model.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dbhelper.UpdateTodoRequest(todoID, userID, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "todo updated successfully",
	})
}

func DeleteTodo(c *gin.Context) {
	auth, ok := util.GetAuth(c)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userID := auth.UserID
	todoID := c.Param("id")

	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "todo id is required"})
		return
	}

	err := dbhelper.DeleteTodo(todoID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "todo deleted successfully"})
}

func GetTodoById(c *gin.Context) {
	auth, ok := util.GetAuth(c)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userID := auth.UserID
	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "todo id is required"})
		return
	}

	todo, err := dbhelper.GetTodoByID(userID, todoID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "todo not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func GetTodos(c *gin.Context) {

	auth, ok := util.GetAuth(c)
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userID := auth.UserID
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

	todos, err := dbhelper.GetTodosByStatus(userID, status, search, limit, offset)
	if err != nil {

		fmt.Println("er")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}
