package server

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
)

func SendChunksToPlayer(p *game.Player, chunks []*game.Chunk) {
	for _, c := range chunks {
		chunkPacket := []byte{packets.PacketChunk}
		chunkPacket = append(chunkPacket, c.ToBytes()...)
		chunkPacket = append(chunkPacket, '\n')
		p.Conn.Write(chunkPacket)
	}
}
