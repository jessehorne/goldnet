package components

import (
	"github.com/gdamore/tcell/v2"
)

type SpriteComponent struct {
	Character  rune
	Foreground tcell.Color
	Background tcell.Color
}

func NewSpriteComponent(character rune, foreground, background tcell.Color) *SpriteComponent {
	return &SpriteComponent{
		Character:  character,
		Foreground: foreground,
		Background: background,
	}
}
