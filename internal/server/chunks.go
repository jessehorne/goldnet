package server

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
)

func SendChunksToPlayer(p *game.Player, chunks []*game.Chunk) {
	for _, c := range chunks {
		p.Conn.Write(packets.BuildChunkPacket(c.ToBytes()))
	}
}
