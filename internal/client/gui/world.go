package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/rivo/tview"
	"sync"
)

type World struct {
	Root    *tview.Box
	Focused bool
	Chunks  []*game.Chunk
	OffsetX int
	OffsetY int
	Mutex   sync.Mutex
}

func NewWorld() *World {
	box := tview.NewBox().SetBorder(false)
	m := &World{
		Root:    box,
		Focused: true,
		Chunks:  []*game.Chunk{},
		OffsetX: 60,
		OffsetY: 10,
	}
	box = box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		m.Draw(screen, x, y, width, height)
		return x, y, width, height
	})
	return m
}

func (m *World) Draw(screen tcell.Screen, x, y, width, height int) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	for _, c := range m.Chunks {
		startX := c.X * game.CHUNK_W
		startY := c.Y * game.CHUNK_H
		for yy := int64(0); yy < game.CHUNK_H; yy++ {
			for xx := int64(0); xx < game.CHUNK_W; xx++ {
				bx := m.OffsetX + int(startX+xx)
				by := m.OffsetY + int(startY+yy)
				if bx > 0 && bx < width && by > 0 && by < 26 {
					screen.SetContent(bx, by, rune(c.Data[yy][xx]), nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
				}
			}
		}
	}

	screen.SetContent(width/2, height/2, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
}

func (m *World) UpdateChunks(chunks []*game.Chunk) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Chunks = chunks
}
