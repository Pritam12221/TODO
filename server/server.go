package server

import (
	handlers "TODO/handlers"
	"TODO/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServerRoutes() *gin.Engine {

	r := gin.Default()

	serverCheck := r.Group("/v1")
	{
		serverCheck.POST("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "server chalrha",
			})
		})
	}
	//user routes
	userRoutes := r.Group("/v1")
	{
		userRoutes.POST("/register", handlers.RegisterUser)
		userRoutes.POST("/login", handlers.LoginUser)
		userRoutes.POST("/renew", handlers.RenewToken)
	}

	//admin auth routes
	adminRoutes := r.Group("/v1/admin")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminRoutes.PATCH("/:id/suspend", handlers.SuspendUser)
		adminRoutes.PATCH("/:id/unsuspend", handlers.UnsuspendUser)
		adminRoutes.GET("/todos", handlers.GetAllTodos)
		adminRoutes.GET("/users", handlers.FetchAllUsers)
	}

	//user auth routes
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
	return r
}
