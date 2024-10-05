package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/server/packets"
	"net"
)

func ClientUserJoinHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
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
