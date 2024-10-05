package game

import (
	"github.com/jessehorne/goldnet/internal/server/packets"
	"github.com/jessehorne/goldnet/internal/util"
	"log"
	"net"
	"os"
)

type GameState struct {
	Players map[int64]*Player
	Logger  *log.Logger
	Chunks  map[int64]map[int64]*Chunk
}

func NewGameState() *GameState {
	return &GameState{
		Logger:  log.New(os.Stdout, "[GoldNet] (GameState) ", log.Ldate|log.Ltime),
		Players: map[int64]*Player{},
		Chunks:  map[int64]map[int64]*Chunk{},
	}
}

func (gs *GameState) GetPlayer(playerID int64) *Player {
	p, ok := gs.Players[playerID]
	if !ok {
		return nil
	}
	return p
}

func (gs *GameState) AddPlayer(playerID int64, c net.Conn) {
	gs.Players[playerID] = NewPlayer(playerID, 0, 0, c)
	gs.SendNearbyChunksToPlayer(gs.Players[playerID])
}

func (gs *GameState) RemovePlayer(playerID int64) {
	gs.Players[playerID].Conn.Close()
	gs.Players[playerID] = nil
}

func (gs *GameState) HandlePlayerAction(playerID int64, action byte) {
	p := gs.GetPlayer(playerID)
	if p != nil {
		p.Action(action)

		if packets.IsMovementAction(action) {
			gs.HandlePlayerMovementAction(p, action)
		}
	}
}

func (gs *GameState) HandlePlayerMovementAction(p *Player, action byte) {
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
	newChunkX := p.X / CHUNK_W
	newChunkY := p.Y / CHUNK_H
	if newChunkX != p.OldChunkX || newChunkY != p.OldChunkY {
		p.OldChunkX = p.X / CHUNK_W
		p.OldChunkY = p.Y / CHUNK_H
		gs.SendNearbyChunksToPlayer(p)
	}
}

func (gs *GameState) SendNearbyChunksToPlayer(p *Player) {
	chunksToSend := gs.GetChunksAroundPlayer(p)
	for _, c := range chunksToSend {
		chunkPacket := []byte{packets.PacketChunk}
		chunkPacket = append(chunkPacket, c.ToBytes()...)
		chunkPacket = append(chunkPacket, '\n')
		p.Conn.Write(chunkPacket)
	}
}

func (gs *GameState) GetChunksAroundPlayer(p *Player) []*Chunk {
	var chunks []*Chunk
	for y := p.OldChunkY - p.ChunkDistance; y < p.OldChunkY+p.ChunkDistance; y++ {
		for x := p.OldChunkX - p.ChunkDistance; x < p.OldChunkX+p.ChunkDistance; x++ {
			_, ok := gs.Chunks[y][x]
			if !ok {
				newChunk := NewChunk(x, y)
				newChunk.FillV1()
				_, yExists := gs.Chunks[y]
				if !yExists {
					gs.Chunks[y] = map[int64]*Chunk{}
				}
				gs.Chunks[y][x] = newChunk
			}
			chunks = append(chunks, gs.Chunks[y][x])
		}
	}
	return chunks
}
