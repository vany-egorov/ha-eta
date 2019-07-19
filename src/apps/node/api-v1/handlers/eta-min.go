package h

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ETAMin(c *gin.Context) {
	name, _ := os.Hostname()

	c.JSON(200, gin.H{
		"hostname": name,
	})
}
