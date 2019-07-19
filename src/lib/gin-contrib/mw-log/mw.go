package mwLog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	logLevel "github.com/vany-egorov/ha-eta/lib/log-level"
)

func New(fnLog func(logLevel.Level)) gin.HandlerFunc {
	return NewWithFilter(fnLog, func(c *gin.Context) bool { return true })
}

func NewWithFilter(fnLog func(logLevel.Level), fnFilter func(*gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if fn := fnFilter; fn != nil {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		end := time.Now()
		httpStatus := c.Writer.Status()
		prefix := c.MustGet("prefix").(string)

		// http server get polled by other clients
		// to prevent spam logging downgrading
		// log-level from info to debug
		isPolling := false
		if v, ok := c.Get("is-polling"); ok {
			if vb, ok := v.(bool); ok {
				isPolling = vb
			}
		}

		buf := bufPool.NewBuf()
		buf.WriteString(prefix)
		buf.WriteString(" | ")
		buf.WriteString(fmt.Sprintf("%3d", httpStatus))
		buf.WriteString(" | ")
		buf.WriteString(fmt.Sprintf("%13v", end.Sub(start)))
		buf.WriteString(" | ")
		buf.WriteString(fmt.Sprintf("%13s", c.ClientIP()))
		buf.WriteString(" | ")
		buf.WriteString(fmt.Sprintf("%4s", c.Request.Method))
		buf.WriteString(" ")
		buf.WriteString(fmt.Sprintf("%22s", c.Request.URL.RequestURI()))
		buf.WriteString(" | ")
		buf.WriteString(fmt.Sprintf("in -> %4d %s", c.Request.ContentLength, c.ContentType()))
		buf.WriteString(" | ")
		buf.WriteString(fmt.Sprintf("out <- %4d %s", c.Writer.Size(), c.Writer.Header().Get("Content-Type")))
		buf.WriteString(" | ")
		buf.WriteString(c.Request.UserAgent())

		switch {
		case httpStatus >= http.StatusInternalServerError:
			fnLog(logLevel.Error, w.String())
		case httpStatus >= http.StatusBadRequest:
			fnLog(logLevel.Warn, w.String())
		case isPolling:
			fnLog(logLevel.Debug, w.String())
		default:
			fnLog(logLevel.Info, w.String())
		}

		for _, e := range c.Errors {
			logger.GetLogger().Errorf("%s | %s", prefix, e.Error())
		}

		pool.Put(w)
	}
}
