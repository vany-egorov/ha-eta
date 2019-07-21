package models

type Point struct {
	Lat float64 `form:"lat" json:"lat" binding:"required"`
	Lng float64 `form:"lng" json:"lng" binding:"required"`
}
