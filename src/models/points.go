package models

import (
	"github.com/vany-egorov/ha-eta/lib/geo-engine/wheely"
)

type Points []Point

/* impl From<wheely.Cars> for *Points */
func (it *Points) FromWheelyCars(cars wheely.Cars) {
	for _, car := range cars {
		*it = append(*it, Point{Lat: car.Lat, Lng: car.Lng})
	}
}

/* impl To<wheely.Points> for Points */
func (it Points) ToWheelyPoints(tgt *wheely.Points) {
	for _, p := range it {
		tgt.Push(p.Lat, p.Lng)
	}
}
