package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/shared"
	packets "github.com/jessehorne/goldnet/packets/dist"

	"github.com/jessehorne/goldnet/internal/game/components"

	"github.com/jessehorne/goldnet/internal/game"
)

func ServerUserJoinHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user joined with a ID of", playerID)

	// add player to gamestates list of players
	newPlayer := components.NewPlayer(components.EntityId(playerID), nil, conn)
	gs.InitNewPlayer(newPlayer)
}

func ServerUserDisconnectedHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user disconnected")

	// remove player from gamestate
	gs.RemovePlayer(components.EntityId(playerID))

	// let everyone know they left
	game.SendOneToAll(gs, &packets.PlayerDisconnected{
		Type: shared.PacketPlayerDisconnected,
		Id:   playerID,
	})
}
