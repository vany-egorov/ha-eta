package seelog

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func New(logger LoggerGetter) gin.HandlerFunc {
	return NewCommmon(logger, func(c *gin.Context) bool { return true })
}

func NewCommmon(logger LoggerGetter, filter func(*gin.Context) bool) gin.HandlerFunc {
	pool := &sync.Pool{
		New: func() interface{} {
			buf := new(bytes.Buffer)
			return buf
		},
	}

	return func(c *gin.Context) {
		if !filter(c) {
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

		w := pool.Get().(*bytes.Buffer)
		w.Reset()
		w.WriteString(prefix)
		w.WriteString(" | ")
		w.WriteString(fmt.Sprintf("%3d", httpStatus))
		w.WriteString(" | ")
		w.WriteString(fmt.Sprintf("%13v", end.Sub(start)))
		w.WriteString(" | ")
		w.WriteString(fmt.Sprintf("%13s", c.ClientIP()))
		w.WriteString(" | ")
		w.WriteString(fmt.Sprintf("%4s", c.Request.Method))
		w.WriteString(" ")
		w.WriteString(fmt.Sprintf("%22s", c.Request.URL.RequestURI()))
		w.WriteString(" | ")
		w.WriteString(fmt.Sprintf("in -> %4d %s", c.Request.ContentLength, c.ContentType()))
		w.WriteString(" | ")
		w.WriteString(fmt.Sprintf("out <- %4d %s", c.Writer.Size(), c.Writer.Header().Get("Content-Type")))
		w.WriteString(" | ")
		w.WriteString(c.Request.UserAgent())

		switch {
		case httpStatus >= http.StatusInternalServerError:
			logger.GetLogger().Error(w.String())
		case httpStatus >= http.StatusBadRequest:
			logger.GetLogger().Warn(w.String())
		case isPolling:
			logger.GetLogger().Debugf(w.String())
		default:
			logger.GetLogger().Info(w.String())
		}

		for _, e := range c.Errors {
			logger.GetLogger().Errorf("%s | %s", prefix, e.Error())
		}

		pool.Put(w)
	}
}
