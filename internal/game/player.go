package game

import "net"

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
