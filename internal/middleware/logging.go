package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] %s %s %d %v",
			time.Now().Format(time.RFC3339),
			method,
			path,
			status,
			latency,
		)
	}
}

// Recovery returns a recovery middleware.
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
