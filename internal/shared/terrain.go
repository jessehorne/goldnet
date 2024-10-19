package shared

import (
	"github.com/gdamore/tcell/v2"
)

const (
	TerrainNothing byte = iota
	TerrainWater
	TerrainSand
	TerrainGrass
)

const (
	BlockNothing byte = iota
	BlockTree
	BlockGrass
	BlockWater
	BlockSand
)

var (
	StyleDefault = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
)

var (
	TerrainStyles = map[byte]tcell.Style{
		TerrainNothing: StyleDefault,
		TerrainWater:   tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue),
		TerrainSand:    tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow),
		TerrainGrass:   tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen),
	}

	Terrain = map[byte]byte{
		BlockNothing: ' ',
		BlockTree:    'T',
		BlockGrass:   ' ',
		BlockWater:   '~',
		BlockSand:    '.',
	}
)

func GetTerrainStyle(height byte) tcell.Style {
	s, ok := TerrainStyles[GetTerrainType(height)]
	if !ok {
		return TerrainStyles[TerrainNothing]
	}
	return s
}

func GetTerrainType(height byte) byte {
	if height < 80 {
		return TerrainWater
	}

	if height >= 80 && height < 100 {
		return TerrainSand
	}

	return TerrainGrass
}
