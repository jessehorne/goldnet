package handlers

import (
	"fmt"
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
)

func ClientMessageHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	playerID, msg := packets.ParseMessagePacket(data)
	var who string
	if playerID == -1 {
		who = "(GAME)"
	} else {
		who = fmt.Sprintf("Player #%d", playerID)
	}
	g.Chat.AddMessage(fmt.Sprintf("%s - %s", who, msg))
}
