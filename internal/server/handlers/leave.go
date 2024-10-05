package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/server/packets"
	"net"
)

func ClientUserLeaveHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user disconnected")

	// remove player from gamestate
	gs.RemovePlayer(playerID)

	// let everyone know they left
	for _, p := range gs.Players {
		if p == nil {
			continue
		}
		p.Conn.Write([]byte{packets.PacketPlayerDisconnected, '\n'})
	}
}
