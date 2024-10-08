package game

import (
	"github.com/jessehorne/goldnet/internal/util"
	"log"
	"os"
	"sync"
)

type GameState struct {
	Players  map[int64]*Player
	Mutex    sync.Mutex
	Logger   *log.Logger
	Chunks   map[int64]map[int64]*Chunk
	IntStore map[string]int64
}

func NewGameState() *GameState {
	return &GameState{
		Logger:   log.New(os.Stdout, "[GoldNet] (GameState) ", log.Ldate|log.Ltime),
		Players:  map[int64]*Player{},
		Chunks:   map[int64]map[int64]*Chunk{},
		IntStore: map[string]int64{},
	}
}

func (gs *GameState) GetPlayer(playerID int64) *Player {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	p, ok := gs.Players[playerID]
	if !ok {
		return nil
	}
	return p
}

func (gs *GameState) AddPlayer(p *Player) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	gs.Players[p.ID] = p
}

func (gs *GameState) RemovePlayer(playerID int64) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	delete(gs.Players, playerID)
}

func (gs *GameState) UpdatePlayerChunks(playerID, x, y int64) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	_, ok := gs.Players[playerID]
	if ok {
		gs.Players[playerID].OldChunkX = x
		gs.Players[playerID].OldChunkY = y
	}
}

func (gs *GameState) SetIntStore(key string, value int64) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	gs.IntStore[key] = value
}

func (gs *GameState) GetIntStore(key string) (int64, bool) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	v, ok := gs.IntStore[key]
	return v, ok
}

func (gs *GameState) MovePlayer(playerID, x, y int64) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	p := gs.Players[playerID]
	if p != nil {
		p.X = x
		p.Y = y
	}
}

func (gs *GameState) HandlePlayerAction(player *Player, action byte) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	if player != nil {
		player.Action(action)
	}
}

func (gs *GameState) GetChunksAroundPlayer(p *Player) []*Chunk {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	var chunks []*Chunk
	for y := p.OldChunkY - 3; y < p.OldChunkY+3; y++ {
		for x := p.OldChunkX - 11; x < p.OldChunkX+10; x++ {
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

func (gs *GameState) GetPlayersAroundPlayer(p *Player) []*Player {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	var players []*Player
	if p == nil {
		return players
	}
	for _, otherPlayer := range gs.Players {
		if otherPlayer == nil {
			continue
		}
		if p.ID == otherPlayer.ID {
			continue
		}

		if util.Distance(p.X, p.Y, otherPlayer.X, otherPlayer.Y) < 10 {
			players = append(players, otherPlayer)
		}
	}
	return players
}

func (gs *GameState) AddChunk(chunk *Chunk) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	_, ok := gs.Chunks[chunk.Y]
	if !ok {
		gs.Chunks[chunk.Y] = map[int64]*Chunk{}
	}
	gs.Chunks[chunk.Y][chunk.X] = chunk
}
