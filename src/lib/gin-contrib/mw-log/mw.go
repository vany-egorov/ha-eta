package mwLog

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	bufPool "github.com/vany-egorov/ha-eta/lib/buf-pool"
	"github.com/vany-egorov/ha-eta/lib/log"
)

type fnLogType func(log.Level, string)
type fnFilterType func(*gin.Context) bool

var toDiscard fnLogType = func(lvl log.Level, msg string) {
	fmt.Fprint(ioutil.Discard, msg)
}

func New(fnLog fnLogType) gin.HandlerFunc {
	return NewWithFilter(fnLog, func(c *gin.Context) bool { return false })
}

func NewWithFilter(fnLog fnLogType, fnFilter fnFilterType) gin.HandlerFunc {
	return func(c *gin.Context) {
		if fn := fnFilter; fn != nil && fn(c) {
			c.Next()
			return
		}

		if fnLog == nil {
			fnLog = toDiscard
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
		defer buf.Release()

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
			fnLog(log.Error, buf.String())
		case httpStatus >= http.StatusBadRequest:
			fnLog(log.Warn, buf.String())
		case isPolling:
			fnLog(log.Debug, buf.String())
		default:
			fnLog(log.Info, buf.String())
		}

		for _, e := range c.Errors {
			buf.Reset()
			buf.WriteString(prefix)
			buf.WriteString(" | ")
			buf.WriteString(e.Error())

			fnLog(log.Error, buf.String())
		}
	}
}
