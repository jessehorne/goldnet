package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"net"
)

func ClientChunksHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var chunksPacket packets.Chunks
	err := proto.Unmarshal(data, &chunksPacket)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal chunks byte data into the chunks packet")
		return
	}

	chunks := game.ParseChunksFromBytes(chunksPacket.Data)
	gs.AddChunks(chunks)

	playerID, ok := gs.GetIntStore("playerID")
	if !ok {
		return
	}
	p := gs.GetPlayer(playerID)
	nearbyChunks, _ := gs.GetChunksAroundPlayer(p)
	g.World.UpdateChunks(nearbyChunks)
}
