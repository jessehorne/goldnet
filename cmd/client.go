package main

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/client"
	"github.com/rivo/tview"
	"os"
	"runtime/debug"
	"time"
)

func main() {
	debug.SetPanicOnFault(true)

	tv := tview.NewApplication()
	c, err := client.NewClient(tv)
	if err != nil {
		tv.Stop()
		panic(err)
	}

	go func() {
		c.Listen()
	}()

	defer func() {
		if e := recover(); e != nil {
			tv.Stop()
			c.Close()
			os.WriteFile("crash.log", []byte(e.(error).Error()), 0644)
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
		tv.Stop()
		panic(err)
	}
	c.Close()
}
