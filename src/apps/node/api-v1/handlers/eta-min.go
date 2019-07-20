package h

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/vany-egorov/ha-eta/apps/node/api-v1"
	apiErrors "github.com/vany-egorov/ha-eta/apps/node/api-v1/errors"
	cache "github.com/vany-egorov/ha-eta/lib/cache"
	geoEngine "github.com/vany-egorov/ha-eta/lib/geo-engine"
)

type ETAMinCtx interface {
	GeoEngine() geoEngine.Engine
	Cache() cache.Cache
}

func ETAMin(c *gin.Context) {
	req := v1.ReqETAMin{}

	if err := c.Bind(&req); err != nil {
		v1.Send(c.Writer, apiErrors.ETAMinReqParse)
		return
	}

	if err := req.Validate(); err != nil {
		v1.Send(c.Writer, err)
		return
	}
}
