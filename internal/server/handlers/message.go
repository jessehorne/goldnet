package handlers

import (
	"net"

	packets2 "github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
)

func ServerMessageHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	msg := packets2.ParseSendMessagePacket(data)
	for _, p := range gs.Players {
		if p == nil {
			return
		}
		p.Conn.Write(packets.BuildMessagePacket(playerID, msg))
	}
}
