package components

import (
	"net"
	"time"

	"github.com/jessehorne/goldnet/internal/game/inventory"
)

type PlayerComponent struct {
	ID        EntityId
	OldChunkX int64
	OldChunkY int64
	Inventory *inventory.Inventory
	Conn      net.Conn

	LastMovementTime time.Time
	LastAttackTime   time.Time

	Username    string
	Gold        int64
	HP          int64
	ST          int64
	Speed       byte    // how many blocks per second the player can travel (water speed is Speed/2)
	AttackSpeed float32 // how many times the player can attack per second

	Hostile bool
}

func NewPlayer(entityId EntityId, inv []byte, c net.Conn) *PlayerComponent {
	return &PlayerComponent{
		ID:               entityId,
		Inventory:        inventory.NewInventory(inv),
		Conn:             c,
		LastMovementTime: time.Now(),
		LastAttackTime:   time.Now(),

		Username:    "bob",
		Gold:        0,
		HP:          10,
		ST:          2,
		Speed:       10,
		AttackSpeed: 1,

		Hostile: false,
	}
}
