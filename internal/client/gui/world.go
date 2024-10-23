package gui

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared"
	"github.com/rivo/tview"
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
		for blockY := int64(0); blockY < game.CHUNK_H; blockY++ {
			for blockX := int64(0); blockX < game.CHUNK_W; blockX++ {
				bx := m.OffsetX + int(startX+blockX)
				by := m.OffsetY + int(startY+blockY)

				// make sure we aren't drawing out of bounds
				if bx > 0 && bx < width && by > 0 && by < 26 {
					top := c.GetTopBlock(blockX, blockY)
					bottom := c.Stack[blockY][blockX][0]
					screen.SetContent(bx, by, rune(top), nil, shared.GetTerrainStyle(bottom))
				}
			}
		}
	}

	// draw sprites
	for entityId, sprite := range m.GameState.SpriteComponents {
		if sprite != nil {
			entityPosition := m.GameState.PositionComponents[entityId]
			if entityPosition != nil {
				bx := m.OffsetX + int(entityPosition.X)
				by := m.OffsetY + int(entityPosition.Y)
				if bx > 0 && bx < width && by > 0 && by < 26 {
					style := tcell.StyleDefault.Foreground(sprite.Foreground).Background(sprite.Background)
					screen.SetContent(bx, by, sprite.Character, nil, style)
				}
			}
		}
	}
}

func (m *World) UpdateChunks(chunks []*game.Chunk) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Chunks = chunks
}
