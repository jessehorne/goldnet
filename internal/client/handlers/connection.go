package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"net"
)

func ClientPlayerJoinedHandler(gs *game.GameState, conn net.Conn, data []byte) {
	gs.Logger.Println("a user joined")
}

func ClientPlayerDisconnectedHandler(gs *game.GameState, conn net.Conn, data []byte) {
	gs.Logger.Println("a user disconnected")
}
