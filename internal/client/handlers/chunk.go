package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"net"
)

func ClientChunkHandler(gs *game.GameState, conn net.Conn, data []byte) {
	chunk := game.ParseChunkFromBytes(data)
	gs.Logger.Println("received chunk", chunk.X, chunk.Y)
}
