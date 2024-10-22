package handlers

import (
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
)

func ClientMessageHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var msg packets.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal message packet")
		return
	}
	g.Chat.AddMessage(msg.Data)
}
