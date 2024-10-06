package server

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

func HandlePlayerMovementAction(gs *game.GameState, p *game.Player, action byte) {
	// send movement to any players nearby
	movePacket := packets.BuildMovePacket(p.ID, p.X, p.Y)
	for _, o := range gs.Players {
		if o != nil {
			if o.ID != p.ID {
				if util.Distance(o.X, o.Y, p.X, p.Y) < 100 {
					p.Conn.Write(movePacket)
				}
			}
		}
	}

	// check if the moving user needs new chunks
	newChunkX := p.X / game.CHUNK_W
	newChunkY := p.Y / game.CHUNK_H
	if newChunkX != p.OldChunkX || newChunkY != p.OldChunkY {
		p.OldChunkX = p.X / game.CHUNK_W
		p.OldChunkY = p.Y / game.CHUNK_H
		nearbyChunks := gs.GetChunksAroundPlayer(p)
		SendChunksToPlayer(p, nearbyChunks)
	}
}
