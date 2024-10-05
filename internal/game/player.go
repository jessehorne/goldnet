package game

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/server/packets"
	"net"
)

type Player struct {
	ID        int64
	X         int64
	Y         int64
	Sprite    rune
	Health    byte
	Inventory []InventoryItem
	Conn      net.Conn
}

func NewPlayer(id, x, y int64, c net.Conn) *Player {
	return &Player{
		ID:        id,
		X:         x,
		Y:         y,
		Sprite:    '@',
		Health:    255,
		Inventory: []InventoryItem{},
		Conn:      c,
	}
}

func (p *Player) Action(a byte) {
	if a == packets.ActionMoveLeft {
		p.X--
		fmt.Println("a player moved left")
	} else if a == packets.ActionMoveRight {
		p.X++
	} else if a == packets.ActionMoveUp {
		p.Y--
	} else if a == packets.ActionMoveDown {
		p.Y++
	}
}
