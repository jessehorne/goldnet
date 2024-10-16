package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/util"
)

func ClientNewZombieHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	c := 0
	id := util.BytesToInt64(data[c : c+8])
	c += 8
	x := util.BytesToInt64(data[c : c+8])
	c += 8
	y := util.BytesToInt64(data[c : c+8])
	c += 8
	hp := util.BytesToInt64(data[c : c+8])
	c += 8
	dmg := util.BytesToInt64(data[c : c+8])
	c += 8
	gold := util.BytesToInt64(data[c : c+8])
	c += 8
	followingPlayerId := util.BytesToInt64(data[c : c+8])
	newZombie := &game.Zombie{
		ID:                id,
		X:                 x,
		Y:                 y,
		HP:                hp,
		Damage:            dmg,
		GoldDropAmt:       gold,
		FollowingPlayerId: followingPlayerId,
	}

	gs.Zombies[id] = newZombie
}

func ClientUpdateZombieHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	c := 0
	id := util.BytesToInt64(data[c : c+8])
	c += 8
	x := util.BytesToInt64(data[c : c+8])
	c += 8
	y := util.BytesToInt64(data[c : c+8])
	c += 8
	hp := util.BytesToInt64(data[c : c+8])
	c += 8
	dmg := util.BytesToInt64(data[c : c+8])
	c += 8
	gold := util.BytesToInt64(data[c : c+8])
	c += 8
	followingPlayerId := util.BytesToInt64(data[c : c+8])
	c += 8

	z, exists := gs.Zombies[id]
	if exists {
		z.X = x
		z.Y = y
		z.HP = hp
		z.Damage = dmg
		z.GoldDropAmt = gold
		z.FollowingPlayerId = followingPlayerId
	}
}
