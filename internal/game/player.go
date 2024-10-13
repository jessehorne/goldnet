package game

import (
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
	"net"
	"time"
)

type Player struct {
	ID               int64
	X                int64
	Y                int64
	OldChunkX        int64
	OldChunkY        int64
	Sprite           rune
	Health           byte
	Inventory        []InventoryItem
	Conn             net.Conn
	Speed            byte // how many blocks per second the player can travel (water speed is Speed/2)
	LastMovementTime time.Time

	Username string
	HP       int64
	ST       int64
}

func NewPlayer(id, x, y int64, c net.Conn) *Player {
	return &Player{
		ID:               id,
		X:                x,
		Y:                y,
		Sprite:           '@',
		Health:           255,
		Inventory:        []InventoryItem{},
		Conn:             c,
		Speed:            10,
		LastMovementTime: time.Now(),

		HP: 10,
		ST: 2,
	}
}

func (p *Player) ToBytes() []byte {
	var data []byte

	// Start with length of username
	data = append(data, util.Int64ToBytes(int64(len(p.Username)))...)
	// Add username bytes
	data = append(data, []byte(p.Username)...)
	// Add X and Y coordinates
	data = append(data, util.Int64ToBytes(p.X)...)
	data = append(data, util.Int64ToBytes(p.Y)...)
	// add HP and ST
	data = append(data, util.Int64ToBytes(p.HP)...)
	data = append(data, util.Int64ToBytes(p.ST)...)

	return data
}

func (p *Player) Action(a byte) {
	if a == packets.ActionMoveLeft {
		p.X--
	} else if a == packets.ActionMoveRight {
		p.X++
	} else if a == packets.ActionMoveUp {
		p.Y--
	} else if a == packets.ActionMoveDown {
		p.Y++
	}
}
