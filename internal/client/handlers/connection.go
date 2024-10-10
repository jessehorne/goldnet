package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
	"net"
)

func ClientPlayerSelfJoinedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	playerID, x, y, others := packets.ParsePlayerSelfJoinedPacket(data)
	p := &game.Player{
		ID:     playerID,
		X:      x,
		Y:      y,
		Sprite: '@',
	}
	gs.AddPlayer(p)
	gs.SetIntStore("playerID", playerID)

	numOfPlayers := util.BytesToInt64(others[0:8])

	if numOfPlayers > 0 {
		counter := 8
		for i := int64(0); i < numOfPlayers; i++ {
			id := util.BytesToInt64(others[counter : counter+8])
			counter += 8
			px := util.BytesToInt64(others[counter : counter+8])
			counter += 8
			py := util.BytesToInt64(others[counter : counter+8])
			counter += 8

			otherPlayer := &game.Player{
				ID: id,
				X:  px,
				Y:  py,
			}
			gs.AddPlayer(otherPlayer)
		}
	}

	g.Chat.AddMessage(util.NewSystemMessage("GAME", "You've connected to GoldNet Official. Good luck!"))
}

func ClientPlayerJoinedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	g.Chat.AddMessage(util.NewSystemMessage("GAME", "A player has appeared!"))

	playerID, x, y := packets.ParsePlayerJoinedPacket(data)
	newPlayer := &game.Player{
		ID: playerID,
		X:  x,
		Y:  y,
	}
	gs.AddPlayer(newPlayer)
}

func ClientPlayerDisconnectedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	g.Chat.AddMessage(util.NewSystemMessage("GAME", "A player has vanished!"))
	playerID := packets.ParsePlayerDisconnectedPacket(data)
	gs.RemovePlayer(playerID)
}
