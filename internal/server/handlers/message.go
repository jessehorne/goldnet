package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
	"net"
)

func ServerMessageHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	msg := packets.ParseSendMessagePacket(data)
	for _, p := range gs.Players {
		if p == nil {
			return
		}
		p.Conn.Write(packets.BuildMessagePacket(playerID, msg))
	}
}
