package handlers

import (
	"github.com/gin-gonic/gin"
)

func NoMethod(c *gin.Context) {
	c.JSON(405, gin.H{
		"code": 405,

		"message": "method is not allowed",
	})
}
