package gui

import "github.com/gdamore/tcell/v2"

const (
	TerrainNothing byte = iota
	TerrainWater
	TerrainSand
	TerrainGrass
	TerrainDirt
	TerrainStone
	TerrainSnow
)

var (
	Terrain = map[byte]tcell.Style{
		TerrainNothing: tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack),
		TerrainWater:   tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
		TerrainSand:    tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorYellow),
		TerrainGrass:   tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite),
		TerrainDirt:    tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorWhite),
		TerrainStone:   tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite),
		TerrainSnow:    tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite),
	}
)

func GetTerrainBlock(d byte) tcell.Style {
	var b byte
	if d < 80 {
		b = TerrainWater
	} else if d < 100 {
		b = TerrainSand
	} else {
		b = TerrainGrass
	}
	_, ok := Terrain[b]
	if !ok {
		return Terrain[0]
	}
	return Terrain[b]
}
