package handlers

import (
	"TODO/database/dbhelper"
	model "TODO/models"
	"TODO/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var req model.UserRequest

	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exists, _:= dbhelper.IsUserExist(req.Email)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error":"User alreaady exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create user
	userID, err := dbhelper.CreateUser(req.Username, req.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store session
	sessionID, sessionErr := dbhelper.CreateUserSession(userID)
	if sessionErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": sessionErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id":    userID,
		"session_id": sessionID,
	})
}

func LoginUser(c *gin.Context) {
	var req model.LoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by mails
	userID, err := dbhelper.GetUserByEmail(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	sessionID, err := dbhelper.CreateUserSession(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":       "User logedin",
		"session_id": sessionID,
	})
}

func LogoutUser(c *gin.Context) {
	
	sessionID := c.GetHeader("session_id")

	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logout failed"})
		return
	}

	err := dbhelper.DeleteUserSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
