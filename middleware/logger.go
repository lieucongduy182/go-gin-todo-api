package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		userID, exist := c.Get("userID")
		userIDStr := "guest"
		if exist {
			userIDStr = string(rune(userID.(uint)))
		}

		// Log format: [TIME] STATUS | LATENCY | IP | USER | METHOD PATH
		log.Printf("[%s] %d | %v | %s | user:%s | %s %s",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			clientIP,
			userIDStr,
			method,
			path,
		)

		if len(c.Errors) > 0 {
			log.Printf("Errors: %v", c.Errors.String())
		}
	}
}
