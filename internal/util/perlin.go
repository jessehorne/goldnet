package util

import (
	"fmt"
	"github.com/aquilax/go-perlin"
)

var (
	perlinGen *perlin.Perlin
)

func PerlinInit(seed int64) {
	perlinGen = perlin.NewPerlin(0.9, 1.9, 4, seed)
}

func PerlinGetByteAtCoords(x, y int64) byte {
	rn := perlinGen.Noise2D(float64(x)/400, float64(y)/200) + 1.0
	formatted := rn * (255 / 2)
	b := byte(formatted)
	if b > 80 && b < 100 {
		fmt.Println(b)
	}
	return b
}
