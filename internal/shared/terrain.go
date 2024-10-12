package shared

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/game"
)

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
		TerrainGrass:   tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorDarkGreen),
		TerrainDirt:    tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorWhite),
		TerrainStone:   tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite),
		TerrainSnow:    tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite),
	}

	Above = map[byte]rune{
		game.AboveNothing: ' ',
		game.AboveTree:    '⤉',
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

func GetTerrainAbove(d byte) rune {
	_, ok := Above[d]
	if ok {
		return Above[d]
	}
	return Above[game.AboveNothing]
}

func GetTerrainBelow(d byte) byte {
	var b byte
	if d < 80 {
		b = TerrainWater
	} else if d < 100 {
		b = TerrainSand
	} else {
		b = TerrainGrass
	}
	return b
}
