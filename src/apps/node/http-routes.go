package node

import "github.com/gin-gonic/gin"

func (it *App) NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	logMw := mwLog.New(func(_) {})

	r := gin.New()
	r.Use(
		gin.Recovery(),
		prefix.New(),
	)

	r.Use(func(c *gin.Context) { c.Set("app-cxt", nil); c.Next() })

	r.HandleMethodNotAllowed = true
	r.NoRoute(handlers.NoRoute, logMw)
	r.NoMethod(handlers.NoMethod, logMw)

	{
		α := r.Group("/api/v1")
		α.Use(logMw)

		α.GET("/eta/min", apiV1.ETAMin)
	}

	return r
}
