package wheely

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Points []Point

func (it *Points) Push(lat float64, lng float64) {
	*it = append(*it, Point{lat, lng})
}

type toPointser interface {
	ToWheelyPoints(*Points)
}

func (it *Points) mustFrom(any interface{}) {
	if that, ok := any.(toPointser); ok {
		that.ToWheelyPoints(it)
	}
}

type Car struct {
	ID    uint64 `json:"id"`
	Point `json:",inline"`
}

type Cars []Car

type fromCarser interface {
	FromWheelyCars(Cars)
}

func (it Cars) mustTo(any interface{}) {
	if that, ok := any.(fromCarser); ok {
		that.FromWheelyCars(it)
	}
}

type PredictReq struct {
	Target Point  `json:"target"`
	Source Points `json:"source"`
}

type ETAs []uint64

type fromETAser interface {
	FromWheelyETAs(ETAs)
}

func (it ETAs) mustTo(any interface{}) {
	if that, ok := any.(fromETAser); ok {
		that.FromWheelyETAs(it)
	}
}
