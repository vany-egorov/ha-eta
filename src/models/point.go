package models

import "fmt"

type Point struct {
	Lat float64 `form:"lat" json:"lat" binding:"required"`
	Lng float64 `form:"lng" json:"lng" binding:"required"`
}

func (it *Point) LatS() string {
	// return strconv.FormatFloat(it.Lat, 'f', -1, 64)
	return fmt.Sprintf("%.7f", it.Lat)
}

func (it *Point) LngS() string {
	// return strconv.FormatFloat(it.Lng, 'f', -1, 64)
	return fmt.Sprintf("%.7f", it.Lng)
}
