package main

import (
	db "TODO/database"
	handlers "TODO/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	err := db.ConnectAndMigrate("localhost", "5433", "mercury-dev", "local", "local", db.SSLModeDisable)

	if err != nil {
		log.Fatal(err)
	}

	v1Group := r.Group("/v1")
	{
		v1Group.POST("/health", func(c *gin.Context) { 	c.JSON(http.StatusOK, gin.H{
			"status":"server working",
	}) })
		v1Group.POST("/register", handlers.RegisterUser)
		v1Group.POST("/login", handlers.LoginUser)
		v1Group.POST("/logout", handlers.LogoutUser)
		v1Group.POST("/todo", handlers.CreateTodo)
		v1Group.PUT("/todo/:id", handlers.UpdateTodo)
  		v1Group.DELETE("/todo/:id", handlers.DeleteTodo)
		v1Group.GET("/todo/:id",handlers.GetTodoById)
		v1Group.GET("/todos",handlers.GetTodos)
	}

	r.Run()
}
