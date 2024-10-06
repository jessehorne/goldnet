package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"net"
)

func ServerActionHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	action := data[0]
	gs.HandlePlayerAction(playerID, action)

	if packets.IsMovementAction(action) {
		// send movement to other players

	}
}
