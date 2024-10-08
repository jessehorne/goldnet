package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"net"
)

func ServerActionHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	action := data[0]

	p := gs.GetPlayer(playerID)
	if p == nil {
		return
	}
	gs.HandlePlayerAction(p, action)

	if packets.IsMovementAction(action) {
		// send movement to other nearby players
		nearbyPlayers := gs.GetPlayersAroundPlayer(p)
		for _, other := range nearbyPlayers {
			other.Conn.Write(packets.BuildMovePacket(p.ID, p.X, p.Y))
		}

		// send chunks if players chunk has updated
		newChunkX := p.X / game.CHUNK_W
		newChunkY := p.Y / game.CHUNK_H
		if newChunkX != p.OldChunkX || newChunkY != p.OldChunkY {
			p.OldChunkX = p.X / game.CHUNK_W
			p.OldChunkY = p.Y / game.CHUNK_H
			nearbyChunks := gs.GetChunksAroundPlayer(p)

			var chunkData []byte
			for _, c := range nearbyChunks {
				chunkData = append(chunkData, c.ToBytes()...)
			}
			chunkPacket := packets.BuildChunksPacket(int64(len(nearbyChunks)), chunkData)
			conn.Write(chunkPacket)
		}
	}
}
