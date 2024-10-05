package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"net"
)

func ClientActionHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	action := data[0]
	gs.HandlePlayerAction(playerID, action)
}
