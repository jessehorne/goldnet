package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"net"
)

func ClientChunksHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	chunks := game.ParseChunksFromBytes(data)
	gs.AddChunks(chunks)
	playerID, ok := gs.GetIntStore("playerID")
	if !ok {
		return
	}
	p := gs.GetPlayer(playerID)
	nearbyChunks, _ := gs.GetChunksAroundPlayer(p)
	g.World.UpdateChunks(nearbyChunks)
}
