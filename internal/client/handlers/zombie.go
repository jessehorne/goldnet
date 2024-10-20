package handlers

import (
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
)

func ClientUpdateZombieHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var z packets.UpdateZombie
	err := proto.Unmarshal(data, &z)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal update zombie packet")
		return
	}
	newZombie := &game.Zombie{
		ID:                z.Id,
		X:                 z.X,
		Y:                 z.Y,
		HP:                z.Hp,
		Damage:            z.Damage,
		GoldDropAmt:       z.GoldDrop,
		FollowingPlayerId: z.FollowingPlayerId,
	}
	gs.Zombies[z.Id] = newZombie
}

func ClientRemoveZombieHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var rz packets.RemoveZombie
	err := proto.Unmarshal(data, &rz)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal remove zombie packet")
		return
	}
	delete(gs.Zombies, rz.Id)
}
