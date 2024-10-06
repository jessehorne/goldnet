package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GUI struct {
	Root  tview.Primitive
	World *World
	Chat  *Chat
	Input *Input
}

func NewGUI() *GUI {
	gui := &GUI{}

	world := NewWorld()
	chat := NewChat()
	input := NewInput()

	grid := tview.NewGrid().
		SetRows(25, 10, 3).
		SetColumns(80).
		SetBorders(true)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		gui.HandleInput(event)
		return event
	})

	grid.AddItem(world.Root, 0, 0, 1, 80, 0, 0, false)
	grid.AddItem(chat.Root, 1, 0, 1, 80, 0, 0, false)
	grid.AddItem(input.Root, 2, 0, 1, 80, 0, 0, false)

	gui.World = world
	gui.Root = grid
	gui.Chat = chat
	gui.Input = input
	return gui
}

func (g *GUI) HandleInput(event *tcell.EventKey) {
	g.Chat.AddMessage(event.Name())
}
