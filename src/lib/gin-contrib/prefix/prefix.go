package prefix

import (
	"github.com/gin-gonic/gin"

	randStr "github.com/vany-egorov/ha-eta/lib/rand-str"
)

func New() gin.HandlerFunc { return NewPrefix() }

func NewPrefix() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("prefix", randStr.Gen(15))
		c.Next()
	}
}
