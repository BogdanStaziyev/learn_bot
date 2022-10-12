package app

import (
	"math"
)

func AtaNumber(x1, y1, x2, y2 float64) (int, int, int) {
	const radius float64 = 180
	const degree, minutes, seconds int = 180, 60, 60
	var deg, min, sec int
	x := x2 - x1
	y := y2 - y1
	subtractionCoordinate := y / x
	atanResult := math.Atan(subtractionCoordinate)
	atanResult *= radius / math.Pi
	deg = int(atanResult)
	minute := (atanResult - float64(deg)) * 60
	min = int(minute)
	sec = int((minute - float64(min)) * 60)
	if x < 0 && y > 0 {
		deg = (degree - 1) + deg
		min = (minutes - 1) + min
		sec = seconds + sec
	} else if x < 0 && y < 0 {
		deg = degree + deg
	} else if x > 0 && y < 0 {
		deg = ((degree * 2) - 1) + deg
		min = (minutes - 1) + min
		sec = seconds + sec
	}
	return deg, min, sec
}
