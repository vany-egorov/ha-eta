package seelog

import (
	"bytes"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func NewApache(logger LoggerGetter) gin.HandlerFunc {
	pool := &sync.Pool{
		New: func() interface{} {
			buf := new(bytes.Buffer)
			return buf
		},
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		c.Next()

		w := pool.Get().(*bytes.Buffer)
		w.Reset()
		w.WriteString(c.ClientIP())
		w.WriteString(" ")
		w.WriteString(time.Now().Format("[02/Jan/2006:15:04:05 -0700] "))
		w.WriteString("\"")
		w.WriteString(c.Request.Method)
		w.WriteString(" ")
		w.WriteString(path)
		w.WriteString(" ")
		w.WriteString(c.Request.Proto)
		w.WriteString("\" ")
		w.WriteString(strconv.Itoa(c.Writer.Status()))
		w.WriteString(" ")
		w.WriteString(strconv.Itoa(c.Writer.Size()))

		logger.GetLogger().Info(w.String())
		pool.Put(w)
	}
}
