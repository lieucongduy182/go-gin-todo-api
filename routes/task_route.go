package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lieucongduy182/go-gin-todo-api/handlers"
	"github.com/lieucongduy182/go-gin-todo-api/middleware"
)

func SetupTaskRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	protected := v1.Group("/tasks")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/", handlers.GetTasks)
		protected.GET("/:id", handlers.GetTask)
		protected.POST("/", handlers.CreateTask)
		protected.PATCH("/:id", handlers.UpdateTask)
		protected.DELETE("/:id", handlers.DeleteTask)
	}
}
