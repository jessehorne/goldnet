package handlers

import (
	"net"

	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/game/components"
)

func ClientUpdateZombieHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var z packets.UpdateZombie
	err := proto.Unmarshal(data, &z)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal update zombie packet")
		return
	}
	newZombie := &components.ZombieComponent{
		ID:                components.EntityId(z.Id),
		HP:                z.Hp,
		Damage:            z.Damage,
		GoldDropAmt:       z.GoldDrop,
		FollowingPlayerId: components.EntityId(z.FollowingPlayerId),
	}
	gs.ZombieComponents[components.EntityId(z.Id)] = newZombie
}

func ClientRemoveZombieHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var rz packets.RemoveZombie
	err := proto.Unmarshal(data, &rz)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal remove zombie packet")
		return
	}
	delete(gs.ZombieComponents, components.EntityId(rz.Id))
}
