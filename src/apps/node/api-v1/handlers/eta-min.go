package h

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/vany-egorov/ha-eta/apps/node/api-v1"
	apiErrors "github.com/vany-egorov/ha-eta/apps/node/api-v1/errors"
)

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
