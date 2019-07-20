package v1

type ReqETAMin struct {
	Lat float64 `form:"lat" json:"lat" binding:"required"`
	Lng float64 `form:"lng" json:"lng" binding:"required"`
}

func (it *ReqETAMin) Validate() error {
	// all validation will be done
	// on upstreams
	return nil
}
