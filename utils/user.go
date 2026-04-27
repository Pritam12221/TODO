package utils

import (
	model "TODO/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

var jwtKey []byte

func InitJWT() {
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
}

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

func GenerateToken(userID, role, sessionID string) (string, error) {

	claims := model.Claims{
		UserID:    userID,
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (*model.Claims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func ParseTokenAllowExpired(tokenStr string) (*model.Claims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if claims, ok := token.Claims.(*model.Claims); ok {
		return claims, nil
	}

	return nil, err
}

func GetUserFromContext(c *gin.Context) (model.User, bool) {
	u, exists := c.Get("user")
	if !exists {
		return model.User{}, false
	}
	user, ok := u.(model.User)
	return user, ok
}
