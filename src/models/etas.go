package models

import (
	"github.com/vany-egorov/ha-eta/lib/geo-engine/wheely"
)

type ETAs []ETA

func (it ETAs) Min() (m ETA) {
	for i, r := range it {
		if i == 0 || r < m {
			m = r
		}
	}

	return
}

/* impl From<wheely.ETAs> for *ETAs */
func (it *ETAs) FromWheelyETAs(etas wheely.ETAs) {
	for _, d := range etas {
		*it = append(*it, ETA(d))
	}
}
