package goCache

import (
	"bytes"
	"strconv"

	"github.com/vany-egorov/ha-eta/models"
)

func keyPoints(w *bytes.Buffer, point models.Point, limit uint64) {
	w.WriteByte('p')
	w.WriteByte('~')
	w.WriteString(point.LatS())
	w.WriteByte(',')
	w.WriteString(point.LngS())
	w.WriteByte('~')
	w.WriteString(strconv.FormatUint(limit, 10))
}

func keyETAs(w *bytes.Buffer, a, b models.Point) {
	w.WriteByte('e')
	w.WriteByte('~')
	w.WriteString(a.LatS())
	w.WriteByte(',')
	w.WriteString(a.LngS())
	w.WriteByte('-')
	w.WriteString(b.LatS())
	w.WriteByte(',')
	w.WriteString(b.LngS())
}
