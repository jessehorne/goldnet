package util

import (
	"math"
	"math/rand"
)

func Distance(x1, y1, x2, y2 int64) int64 {
	xD := float64((x2 - x1) ^ 2)
	yD := float64((y2 - y1) ^ 2)
	return int64(math.Sqrt(xD + yD))
}

func RandomIntBetween(min, max int) int {
	return rand.Intn(max-min) + min
}
