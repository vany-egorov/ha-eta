package geoEngine

type Backernd interface {
	DoCars(lat, lng float64)
	DoPredict() error
}
