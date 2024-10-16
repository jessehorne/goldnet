package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
)

func ClientPlayerToggleHostileHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	id, hostile := packets.ParseSetHostilePacket(data)

	gs.Mutex.Lock()
	gs.Players[id].Hostile = hostile
	gs.Mutex.Unlock()
}
