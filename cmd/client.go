package main

import (
	"github.com/jessehorne/goldnet/internal/client"
	"github.com/rivo/tview"
)

func main() {

	c, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	go func() {
		c.Listen()
	}()

	tv := tview.NewApplication()
	defer func() {
		if err := recover(); err != nil {
			c.Close()
			tv.Stop()
		}
	}()

	if err = tv.SetRoot(c.GUI.Root, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
