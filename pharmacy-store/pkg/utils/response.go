package utils

import (
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, statusCode int, message interface{}) {
	c.JSON(statusCode, gin.H{
		"message": message,
	})
}
