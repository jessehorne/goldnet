package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/rivo/tview"
	"sync"
)

type World struct {
	Root       *tview.Box
	Focused    bool
	Chunks     []*game.Chunk
	OffsetX    int
	OffsetY    int
	OldOffsetX int
	OldOffsetY int
	Mutex      sync.Mutex
	GameState  *game.GameState
}

func NewWorld(gs *game.GameState) *World {
	box := tview.NewBox().SetBorder(false)
	m := &World{
		Root:      box,
		Focused:   true,
		Chunks:    []*game.Chunk{},
		OffsetX:   50,
		OffsetY:   13,
		GameState: gs,
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

	m.GameState.Mutex.Lock()
	defer m.GameState.Mutex.Unlock()

	for _, c := range m.Chunks {
		startX := c.X * game.CHUNK_W
		startY := c.Y * game.CHUNK_H
		for yy := int64(0); yy < game.CHUNK_H; yy++ {
			for xx := int64(0); xx < game.CHUNK_W; xx++ {
				bx := m.OffsetX + int(startX+xx)
				by := m.OffsetY + int(startY+yy)
				if bx > 0 && bx < width && by > 0 && by < 26 {
					b := c.Below[yy][xx]
					screen.SetContent(bx, by, GetTerrainAbove(c.Above[yy][xx]), nil, GetTerrainBlock(b))
				}
			}
		}
	}

	// draw players
	for _, p := range m.GameState.Players {
		if p != nil {
			bx := m.OffsetX + int(p.X)
			by := m.OffsetY + int(p.Y)
			if bx > 0 && bx < width && by > 0 && by < 26 {
				screen.SetContent(bx, by, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}
		}
	}
}

func (m *World) UpdateChunks(chunks []*game.Chunk) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Chunks = chunks
}
