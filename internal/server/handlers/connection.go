package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"net"
)

func ServerUserJoinHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user joined with ID of", playerID)

	// let every player know they joined
	for _, p := range gs.Players {
		if p == nil {
			continue
		}
		p.Conn.Write([]byte{packets.PacketPlayerJoined, '\n'})
	}

	// add player to gamestates list of players
	gs.AddPlayer(playerID, conn)
}

func ServerUserDisconnectedHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
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
