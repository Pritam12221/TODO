package handlers

import (
	"TODO/database/dbhelper"
	model "TODO/models"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	sessionID := c.GetHeader("session_id")

	if sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
		return
	}

	userID, err := dbhelper.GetUserFromSession(sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

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
  sessionID := c.GetHeader("session_id")
  if sessionID == "" {
   c.JSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
   return
  }
  //get userid from session
  userID, err := dbhelper.GetUserFromSession(sessionID)
  if err != nil {
   c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
   return
  }

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

err = dbhelper.UpdateTodoRequest(todoID, userID, req)
if err != nil {
	if err.Error() == "todo not found or already deleted" {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo"})
	return
}

c.JSON(http.StatusOK, gin.H{
	"message": "todo updated successfully",
})}


 func DeleteTodo(c *gin.Context) {
  sessionID := c.GetHeader("session_id")
  if sessionID == "" {
   c.JSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
   return
  }

  userID, err := dbhelper.GetUserFromSession(sessionID)
  if err != nil {
   c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
   return
  }

  todoID := c.Param("id")
  if todoID == "" {
   c.JSON(http.StatusBadRequest, gin.H{"error": "todo id is required"})
   return
  }

  if err := dbhelper.DeleteTodo(todoID, userID); err != nil {
   if err.Error() == "todo not found or not owned by user" {
    c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    return
   }
   c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete todo"})
   return
  }

  c.JSON(http.StatusOK, gin.H{"message": "todo deleted successfully"})
 }


 func GetTodoById(c *gin.Context){
   sessionID := c.GetHeader("session_id")

  if sessionID == "" {
   c.JSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
   return
  }

  userID, err := dbhelper.GetUserFromSession(sessionID)
  if err != nil {
   c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
   return
  }

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
	sessionID := c.GetHeader("session_id")
	if sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
		return
	}

	userID, err := dbhelper.GetUserFromSession(sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	status := c.Query("status")

	todos, err := dbhelper.GetTodosByStatus(userID, status)
	if err != nil {
		if err.Error() == "invalid status" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}