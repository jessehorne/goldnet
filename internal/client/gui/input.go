package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Input struct {
	Root    *tview.Form
	Focused bool
	Value   string
}

func NewInput() *Input {
	newInput := &Input{}

	f := tview.NewForm()
	f.AddInputField("> ", "", 80, nil, func(text string) {
		newInput.Value = text
	})
	f.SetBorderPadding(0, 0, 0, 0)
	f.SetItemPadding(0)
	f.SetBackgroundColor(tcell.ColorBlack)

	newInput.Root = f
	return newInput
}
