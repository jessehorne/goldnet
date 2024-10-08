package main

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/client"
	"github.com/rivo/tview"
	"os"
	"time"
)

func main() {

	tv := tview.NewApplication()
	c, err := client.NewClient(tv)
	if err != nil {
		panic(err)
	}
	go func() {
		c.Listen()
	}()

	defer func() {
		if e := recover(); e != nil {
			os.WriteFile("crash.log", []byte(err.Error()), 0644)
			c.Close()
			tv.Stop()
			fmt.Println(e)
			panic(e)
		}
	}()

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			tv.ForceDraw()
		}
	}()

	if err = tv.SetRoot(c.GUI.Root, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
	c.Close()
}
