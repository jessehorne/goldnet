package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"net"
)

func ClientPlayerMovedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	id, x, y := packets.ParseMovePacket(data)
	gs.MovePlayer(id, x, y)
}
