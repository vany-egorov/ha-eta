package handlers

import (
	"github.com/gin-gonic/gin"
)

func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{
		"code": 404,

		"message": "check if requested URL is valid. maybe misspelled? check docs, configs.",
	})
}
