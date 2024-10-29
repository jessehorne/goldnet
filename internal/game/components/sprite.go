package components

import (
	"github.com/gdamore/tcell/v2"
)

type Sprite struct {
	Character  rune
	Foreground tcell.Color
	Background tcell.Color
}

func NewSpriteComponent(character rune, foreground, background tcell.Color) *Sprite {
	return &Sprite{
		Character:  character,
		Foreground: foreground,
		Background: background,
	}
}
