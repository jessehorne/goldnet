package gui

import "github.com/gdamore/tcell/v2"

const (
	TerrainNothing byte = iota
	TerrainWater
	TerrainDirt
	TerrainGrass
	TerrainStone
	TerrainSnow
)

var (
	Terrain = map[byte]tcell.Style{
		TerrainNothing: tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack),
		TerrainWater:   tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
		TerrainDirt:    tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorWhite),
		TerrainGrass:   tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite),
		TerrainStone:   tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite),
		TerrainSnow:    tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite),
	}
)

func GetTerrainBlock(b byte) tcell.Style {
	_, ok := Terrain[b]
	if !ok {
		return Terrain[0]
	}
	return Terrain[b]
}
