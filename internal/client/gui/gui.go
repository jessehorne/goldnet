package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/rivo/tview"
)

type GUI struct {
	Root    tview.Primitive
	World   *World
	Chat    *Chat
	Input   *Input
	Sidebar *Sidebar
}

func NewGUI(gs *game.GameState, inputFunc func(event *tcell.EventKey) *tcell.EventKey) *GUI {
	gui := &GUI{}

	world := NewWorld(gs)
	chat := NewChat()
	input := NewInput()
	sidebar := NewSidebar()

	grid := tview.NewGrid().
		SetRows(25, 8, 1).
		SetColumns(80).
		SetBorders(true)

	grid.SetInputCapture(inputFunc)

	grid.AddItem(world.Root, 0, 0, 1, 10, 0, 0, true)
	grid.AddItem(sidebar.Root, 0, 10, 1, 20, 0, 0, true)
	grid.AddItem(chat.Root, 1, 0, 1, 30, 0, 0, false)
	grid.AddItem(input.Root, 2, 0, 1, 30, 0, 0, false)

	gui.World = world
	gui.Root = grid
	gui.Chat = chat
	gui.Input = input
	gui.Sidebar = sidebar
	return gui
}

func (g *GUI) HandleInput(event *tcell.EventKey) {
	if g.World.Focused {
		if event.Rune() == 'a' {

		}
	}
}
