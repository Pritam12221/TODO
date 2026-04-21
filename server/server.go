package server

import (
	handlers "TODO/handlers"
	"TODO/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServerRoutes()*gin.Engine{

	r := gin.Default()

	serverCheck :=r.Group("/v1")
	{
	serverCheck.POST("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "server chalrha",
			})
		})
	}
		
	userRoutes := r.Group("/v1")
	{
		userRoutes.POST("/register", handlers.RegisterUser)
		userRoutes.POST("/login", handlers.LoginUser)
	}

	
	userAuth := r.Group("/v1")
	userAuth.Use(middleware.AuthMiddleware())
	{
		userAuth.POST("/logout", handlers.LogoutUser)
		userAuth.POST("/todo", handlers.CreateTodo)
		userAuth.PUT("/todo/:id", handlers.UpdateTodo)
		userAuth.DELETE("/todo/:id", handlers.DeleteTodo)
		userAuth.GET("/todo/:id", handlers.GetTodoById)
		userAuth.GET("/todos", handlers.GetTodos)
	}
	return r;
}