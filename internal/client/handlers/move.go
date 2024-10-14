package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"net"
)

func ClientPlayerMovedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	id, x, y := packets.ParseMovePacket(data)
	playerID, _ := gs.GetIntStore("playerID")
	gs.MovePlayer(id, x, y)
	if playerID == id {
		g.World.OffsetX = 50 + -int(x)
		g.World.OffsetY = 13 + -int(y)
	}
}
