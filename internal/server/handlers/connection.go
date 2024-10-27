package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/shared"
	"github.com/jessehorne/goldnet/internal/util"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"

	"github.com/jessehorne/goldnet/internal/game/components"

	"github.com/jessehorne/goldnet/internal/game"
)

func ServerUserJoinHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user joined with a ID of", playerID)

	// add player to gamestates list of players
	playerId := gs.NextEntityId()
	newPlayer := components.NewPlayer(playerId, nil, conn)
	gs.InitNewPlayer(newPlayer)
}

func ServerUserDisconnectedHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user disconnected")

	// remove player from gamestate
	gs.RemovePlayer(components.EntityId(playerID))

	// let everyone know they left
	dp := &packets.PlayerDisconnected{
		Type: shared.PacketPlayerDisconnected,
		Id:   playerID,
	}

	dpData, dpErr := proto.Marshal(dp)
	if dpErr != nil {
		gs.Logger.Println(dpErr)
		return
	}
	for _, p := range gs.PlayerComponents {
		if p == nil {
			continue
		}
		util.Send(p.Conn, dpData)
	}
}
