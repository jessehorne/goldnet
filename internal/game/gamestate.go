package game

import (
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
}

func (gs *GameState) RemovePlayer(playerID int64) {
	gs.Players[playerID].Conn.Close()
	gs.Players[playerID] = nil
}

func (gs *GameState) HandlePlayerAction(playerID int64, action byte) {
	p := gs.GetPlayer(playerID)
	if p != nil {
		p.Action(action)
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
