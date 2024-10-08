package handlers

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/util"
	"net"
)

func ClientChunksHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	g.Chat.AddMessage(fmt.Sprintf("chunks packet received: %d", len(data)))
	chunks := game.ParseChunksFromBytes(data)
	gs.AddChunks(chunks)

	playerID, ok := gs.GetIntStore("playerID")
	if !ok {
		return
	}
	p := gs.GetPlayer(playerID)
	nearbyChunks := gs.GetChunksAroundPlayer(p)
	g.Chat.AddMessage(fmt.Sprintf("First: %d,%d", nearbyChunks[0].X, nearbyChunks[0].Y))
	g.Chat.AddMessage(fmt.Sprintf("Chunks to Draw: %d", len(nearbyChunks)))
	g.World.UpdateChunks(nearbyChunks)
	//gs.Logger.Println("received chunk", chunk.X, chunk.Y)
	g.Chat.AddMessage(util.NewSystemMessage("GAME", fmt.Sprintf("redrawing %d chunks", len(g.World.Chunks))))
}
