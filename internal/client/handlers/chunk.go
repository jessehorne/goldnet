package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"net"
)

func ClientChunkHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	chunk := game.ParseChunkFromBytes(data)
	gs.AddChunk(chunk)
}
