package handlers

import (
	"github.com/gin-gonic/gin"
)

// TODO: implement swagger handler
func Swagger(c *gin.Context) {
	c.Header("Content-Type", "application/x-yaml")
}
