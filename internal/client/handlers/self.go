package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
)

func ClientUpdatePlayerHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	updatePlayer := packets.ParsePlayerBytes(data)
	p := gs.GetPlayer(updatePlayer.ID)

	p.X = updatePlayer.X
	p.Y = updatePlayer.Y
	p.Gold = updatePlayer.Gold
	p.HP = updatePlayer.HP
	p.ST = updatePlayer.ST
	p.Hostile = updatePlayer.Hostile

	gs.MovePlayer(p.ID, p.X, p.Y)

	currentPlayerID, exists := gs.GetIntStore("playerID")
	if exists {
		if currentPlayerID == updatePlayer.ID {
			g.Sidebar.UpdatePlayerStats(p)
			g.World.OffsetX = 50 + -int(p.X)
			g.World.OffsetY = 13 + -int(p.Y)
		}
	}

}
