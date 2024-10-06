package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type World struct {
	Root tview.Primitive
}

func NewWorld() *World {
	box := tview.NewBox().SetBorder(false)
	m := &World{
		Root: box,
	}
	box = box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		m.Draw(screen, x, y, width, height)
		return x, y, width, height
	})
	return m
}

func (m *World) Draw(screen tcell.Screen, x, y, width, height int) {
	screen.SetContent(width/2, height/2, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
}
