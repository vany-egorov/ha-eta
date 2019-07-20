package handlers

import (
	"github.com/gin-gonic/gin"
)

func NoMethod(c *gin.Context) {
	c.JSON(405, gin.H{
		"message": "method is not allowed",
	})
}
