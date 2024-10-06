package main

import (
	"github.com/jessehorne/goldnet/internal/client"
	"github.com/rivo/tview"
)

func main() {
	gui := client.NewGUI()
	if err := tview.NewApplication().SetRoot(gui.Root, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
