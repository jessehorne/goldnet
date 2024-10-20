package handlers

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/shared"
	"github.com/jessehorne/goldnet/internal/util"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"net"

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

	msgPacket := &packets.Message{
		Type: shared.PacketSendMessage,
		Data: fmt.Sprintf("%s - %s", p.Username, messageFromPlayer.Data),
	}
	msgData, perr := proto.Marshal(msgPacket)
	if perr != nil {
		gs.Logger.Println(perr)
		return
	}
	for _, pl := range gs.Players {
		if pl == nil {
			return
		}
		util.Send(pl.Conn, msgData)
	}
}
