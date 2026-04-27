package utils

import (
	model "TODO/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
}

func GetAuth(c *gin.Context) (model.AuthContext, bool) {
	val, exists := c.Get("auth")
	if !exists {
		return model.AuthContext{}, false
	}

	auth, ok := val.(model.AuthContext)
	return auth, ok
}
