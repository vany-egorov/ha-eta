package models

func FnZipForEarh(points Points, etas ETAs, cb func(Point, ETA)) {
	ln := len(points)
	if lnEtas := len(etas); lnEtas < ln {
		ln = lnEtas
	}

	for i := 0; i < ln; i++ {
		cb(points[i], etas[i])
	}
}
