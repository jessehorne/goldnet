package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
	"net"
)

func ClientPlayerSelfJoinedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	playerID, x, y := packets.ParsePlayerSelfJoinedPacket(data)
	gs.Players[playerID] = &game.Player{
		ID:            playerID,
		X:             x,
		Y:             y,
		ChunkDistance: 4,
		Sprite:        '@',
	}
	gs.IntStore["playerID"] = playerID
	g.Chat.AddMessage(util.NewSystemMessage("GAME", "You've connected to GoldNet Official. Good luck!"))
}

func ClientPlayerJoinedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	g.Chat.AddMessage(util.NewSystemMessage("GAME", "A player has appeared!"))
}

func ClientPlayerDisconnectedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	g.Chat.AddMessage(util.NewSystemMessage("GAME", "A player has vanished!"))
}
