package util

import (
	"github.com/aquilax/go-perlin"
)

var (
	perlinGen *perlin.Perlin
)

const (
	BlockTree byte = iota
)

func PerlinInit(seed int64) {
	perlinGen = perlin.NewPerlin(0.9, 1.9, 4, seed)
}

// PerlinGetDataAtCoords returns what a user will find/see at a specific coordinate
// The first return variable is the type of terrain (water, sand, grass, etc) which can be viewed as height, where
// water is at the lowest level.
func PerlinGetDataAtCoords(x, y int64) byte {
	belowRandom := perlinGen.Noise2D(float64(x)/400, float64(y)/200) + 1.0
	formatted := belowRandom * (255 / 2)
	below := byte(formatted)
	return below
}
