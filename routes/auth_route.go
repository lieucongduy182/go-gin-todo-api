package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lieucongduy182/go-gin-todo-api/handlers"
	"github.com/lieucongduy182/go-gin-todo-api/middleware"
)

func SetupAuthRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
	}

	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", handlers.GetProfile)
	}
}
