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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isExists, _ := dbhelper.IsUserExist(req.Email)
	if isExists {
		c.JSON(http.StatusConflict, gin.H{"error": "User alreaady exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// c.JSON(http.StatusCreated, gin.H{
	// 	"user_id":    userID,
	// 	"session_id": sessionID,
	// })

	token, err := utils.GenerateToken(userID, "employee", sessionID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
		"token":   token,
	})
}

func LoginUser(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := dbhelper.GetUserByEmail(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if user.IsSuspended {
		c.JSON(403, gin.H{"error": "account suspended"})
		return
	}

	sessionID, err := dbhelper.CreateUserSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"user":       "User logedin",
	// 		"session_id": sessionID,
	// 	})

	token, err := utils.GenerateToken(user.ID, string(user.Role), sessionID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

func LogoutUser(c *gin.Context) {
	sessionID := c.GetString("session_id")

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

func RenewToken(c *gin.Context) {

	// Get Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "missing token"})
		return
	}

	// Extract token
	tokenStr := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenStr = authHeader[7:]
	}

	// Parse token even if expired
	claims, err := utils.ParseTokenAllowExpired(tokenStr)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	sessionID := claims.SessionID
	if sessionID == "" {
		c.JSON(401, gin.H{"error": "invalid session in token"})
		return
	}

	// Validate session
	userID, err := dbhelper.GetUserIDBySession(sessionID)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid or expired session"})
		return
	}

	//  Fetch user
	user, err := dbhelper.GetUserByID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "user not found"})
		return
	}

	// Suspension check
	if user.IsSuspended {
		c.JSON(403, gin.H{"error": "account suspended"})
		return
	}

	// Generate new token
	newToken, err := utils.GenerateToken(user.ID, string(user.Role), sessionID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"token":   newToken,
		"message": "new token genrated",
	})
}
