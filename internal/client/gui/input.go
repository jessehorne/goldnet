package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Input struct {
	Root *tview.Form
}

func NewInput() *Input {
	f := tview.NewForm()
	f.AddInputField("> ", "asdsd", 80, nil, nil)
	f.SetBorderPadding(0, 0, 0, 0)
	f.SetItemPadding(0)
	f.SetBackgroundColor(tcell.ColorBlack)

	return &Input{
		Root: f,
	}
}
