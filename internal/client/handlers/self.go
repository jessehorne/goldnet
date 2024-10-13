package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
	"net"
)

func ClientUpdateSelfPlayerHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	g.Chat.AddMessage("test")
	playerID, ok := gs.GetIntStore("playerID")
	if !ok {
		return
	}
	p := gs.GetPlayer(playerID)

	updatePlayer := packets.ParsePlayerBytes(data)
	p.X = updatePlayer.X
	p.Y = updatePlayer.Y
	p.HP = updatePlayer.HP
	p.ST = updatePlayer.ST

	gs.MovePlayer(p.ID, p.X, p.Y)
	g.World.OffsetX = 50 + -int(p.X)
	g.World.OffsetY = 13 + -int(p.Y)

	g.Chat.AddMessage("moved")
}
