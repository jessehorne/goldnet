package handlers

import (
	"fmt"
	"net"

	"github.com/jessehorne/goldnet/internal/shared"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"

	"github.com/jessehorne/goldnet/internal/game"
)

func ServerMessageHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	var messageFromPlayer packets.Message
	err := proto.Unmarshal(data, &messageFromPlayer)
	if err != nil {
		gs.Logger.Println(err)
		return
	}

	p := gs.GetPlayer(playerID)
	if p == nil {
		gs.Logger.Println("failed to send message because player doesn't exist in gamestate")
		return
	}

	game.SendOneToAll(gs, &packets.Message{
		Type: shared.PacketSendMessage,
		Data: fmt.Sprintf("%s - %s", p.Username, messageFromPlayer.Data),
	})
}
