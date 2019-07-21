package v1

import (
	"github.com/vany-egorov/ha-eta/models"
)

type ReqETAMin struct {
	models.Point `json:",inline"`
}

func (it *ReqETAMin) Validate() error {
	// all validation will be done
	// on upstreams
	// just do pass-through
	return nil
}
