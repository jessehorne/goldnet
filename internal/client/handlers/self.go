package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
)

func ClientUpdateSelfPlayerHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	playerID, ok := gs.GetIntStore("playerID")
	if !ok {
		return
	}
	p := gs.GetPlayer(playerID)

	updatePlayer := packets.ParsePlayerBytes(data)
	p.X = updatePlayer.X
	p.Y = updatePlayer.Y
	p.Gold = updatePlayer.Gold
	p.HP = updatePlayer.HP
	p.ST = updatePlayer.ST
	p.Hostile = updatePlayer.Hostile

	gs.MovePlayer(p.ID, p.X, p.Y)
	g.World.OffsetX = 50 + -int(p.X)
	g.World.OffsetY = 13 + -int(p.Y)

	g.Sidebar.UpdatePlayerStats(p)
}
