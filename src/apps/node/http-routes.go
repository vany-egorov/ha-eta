package node

import (
	"github.com/gin-gonic/gin"

	"github.com/vany-egorov/ha-eta/handlers"
	mwLog "github.com/vany-egorov/ha-eta/lib/gin-contrib/mw-log"
	mwPrefix "github.com/vany-egorov/ha-eta/lib/gin-contrib/mw-prefix"
	"github.com/vany-egorov/ha-eta/lib/log"

	apiV1 "github.com/vany-egorov/ha-eta/apps/node/api-v1/handlers"
)

func (it *App) NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// or just `mwLog.New(log.Log)` or `mwLog.New(it.ctx.fnLog)`
	logMw := mwLog.New(func(lvl log.Level, msg string) {
		it.ctx.fnLog(lvl, msg)
	})

	r := gin.New()
	r.Use(
		gin.Recovery(),
		mwPrefix.New(),
	)

	r.HandleMethodNotAllowed = true
	r.NoRoute(handlers.NoRoute, logMw)
	r.NoMethod(handlers.NoMethod, logMw)

	{
		α := r.Group("/api/v1")
		α.Use(logMw)

		α.Use(func(c *gin.Context) { c.Set("service-ctx", &it.ctx); c.Next() })

		α.GET("/eta/min", apiV1.ETAMin)
	}

	return r
}
