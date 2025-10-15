package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/health", healthCheck)

	v1 := router.Group("/api/v1")

	router.Run(":8080")
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running smoothly",
		"time":    time.Now(),
	})
}
