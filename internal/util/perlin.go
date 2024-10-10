package util

import (
	"github.com/aquilax/go-perlin"
)

var (
	perlinGen *perlin.Perlin
)

func PerlinInit(seed int64) {
	perlinGen = perlin.NewPerlin(1.1, 1.5, 6, seed)
}

func PerlinGetByteAtCoords(x, y int64) byte {
	rn := perlinGen.Noise2D(float64(x)/100, float64(y)/50) + 1.0
	formatted := rn * 3
	b := byte(formatted)
	return b
}
