package util

import (
	"github.com/aquilax/go-perlin"
)

var (
	perlinGen *perlin.Perlin
)

func PerlinInit(seed int64) {
	perlinGen = perlin.NewPerlin(1.1, 1.6, 4, seed)
}

func PerlinGetByteAtCoords(x, y int64) byte {
	rn := perlinGen.Noise2D(float64(x)/400, float64(y)/200) + 1.0
	formatted := rn
	b := byte(formatted)
	return b
}
