package handlers

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"net"
)

func ClientPlayerMovedHandler(gs *game.GameState, conn net.Conn, data []byte) {
	id, x, y := packets.ParseMovePacket(data)
	gs.Logger.Println(fmt.Sprintf("player '%d' moved to (%d,%d)", id, x, y))
}
