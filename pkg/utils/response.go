package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSON sends a JSON response with the given status code and data.
func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// Error sends a JSON error response.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// BadRequest sends a 400 Bad Request response.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// InternalError sends a 500 Internal Server Error response.
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
