package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lieucongduy182/go-gin-todo-api/config"
	"github.com/lieucongduy182/go-gin-todo-api/database"
	"github.com/lieucongduy182/go-gin-todo-api/middleware"
	"github.com/lieucongduy182/go-gin-todo-api/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	database.Connect()

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	router.Use(middleware.CustomLogger())

	router.GET("/health", healthCheck)

	// setup routes
	routes.SetupRoutes(router)

	serverAddr := ":" + config.AppConfig.ServerPort
	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running smoothly",
		"time":    time.Now(),
	})
}
