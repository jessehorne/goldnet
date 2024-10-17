package game

import (
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

type GameState struct {
	Players     map[int64]*Player
	Zombies     map[int64]*Zombie
	PlayerCount int64
	Mutex       sync.Mutex
	Logger      *log.Logger
	Chunks      map[int64]map[int64]*Chunk
	IntStore    map[string]int64
	TPS         int // ticks per second
}

func NewGameState() *GameState {
	return &GameState{
		Logger:      log.New(os.Stdout, "[GoldNet] (GameState) ", log.Ldate|log.Ltime),
		Players:     map[int64]*Player{},
		Zombies:     map[int64]*Zombie{},
		PlayerCount: 0,
		Chunks:      map[int64]map[int64]*Chunk{},
		IntStore:    map[string]int64{},
		TPS:         10, // ticks per second
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
		p.OldChunkX = p.X / CHUNK_W
		p.OldChunkY = p.Y / CHUNK_H
	}
}

func (gs *GameState) HandlePlayerAction(player *Player, action byte) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	if player != nil {
		player.Action(action)
	}
}

func (gs *GameState) GetChunkAtCoords(x, y int64) *Chunk {
	cX := x / 8
	cY := y / 8
	c, ok := gs.Chunks[cY][cX]
	if !ok {
		return nil
	}
	return c
}

func (gs *GameState) GetBelowBlockAtCoords(x, y int64) byte {
	c := gs.GetChunkAtCoords(x, y)
	if c == nil {
		return 0
	}
	modX := x % 8
	modY := y % 8
	if modX < 0 {
		modX = 8 - -(modX)
	}
	if modY < 0 {
		modY = 8 - -(modY)
	}
	return c.GetBelowBlock(modX, modY)
}

func (gs *GameState) GetAboveBlockAtCoords(x, y int64) byte {
	c := gs.GetChunkAtCoords(x, y)
	if c == nil {
		return 0
	}
	modX := x - (x / 8)
	modY := y - (y / 8)
	return c.GetAboveBlock(modX, modY)
}

func (gs *GameState) GetChunksAroundPlayer(p *Player) ([]*Chunk, []*Chunk) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	var chunks []*Chunk
	var newChunks []*Chunk
	for y := p.OldChunkY - 3; y < p.OldChunkY+3; y++ {
		for x := p.OldChunkX - 11; x < p.OldChunkX+10; x++ {
			_, ok := gs.Chunks[y][x]
			if !ok {
				newChunk := NewChunk(x, y)
				newChunk.FillPerlin()
				_, yExists := gs.Chunks[y]
				if !yExists {
					gs.Chunks[y] = map[int64]*Chunk{}
				}
				gs.Chunks[y][x] = newChunk
				newChunks = append(newChunks, newChunk)
			}
			chunks = append(chunks, gs.Chunks[y][x])
		}
	}
	return chunks, newChunks
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

		dis := util.Distance(p.X, p.Y, otherPlayer.X, otherPlayer.Y)
		if dis < 50 {
			players = append(players, otherPlayer)
		}
	}
	return players
}

func (gs *GameState) AddChunks(chunks []*Chunk) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	for _, c := range chunks {
		_, ok := gs.Chunks[c.Y]
		if !ok {
			gs.Chunks[c.Y] = map[int64]*Chunk{}
		}
		gs.Chunks[c.Y][c.X] = c
	}
}

func (gs *GameState) UpdateZombies() {

	// update zombie
	gs.Mutex.Lock()
	for _, z := range gs.Zombies {
		// handle movement
		doesMove := util.RandomIntBetween(0, 10) < 2
		if doesMove {
			// If not currently following a player, pick one
			if z.FollowingPlayerId == -1 {
				for _, player := range gs.Players {
					if util.Distance(z.X, z.Y, player.X, player.Y) < ZOMBIE_FOLLOW_RANGE {
						z.FollowingPlayerId = player.ID
					}
				}
			}

			// If we are now following a player, move towards it
			if z.FollowingPlayerId != -1 {
				if z == nil {
					continue
				}
				followingPlayer := gs.Players[z.FollowingPlayerId]
				if followingPlayer == nil {
					continue
				}
				// Follow player if close enough
				if util.Distance(z.X, z.Y, followingPlayer.X, followingPlayer.Y) < ZOMBIE_FOLLOW_RANGE {
					direction := util.RandomIntBetween(0, 2)
					if direction == 0 {
						xDist := followingPlayer.X - z.X
						if xDist*xDist > 0 { // Just checking for positive magnitude
							z.X += xDist / int64(math.Abs(float64(xDist)))
						}
					} else {
						yDist := followingPlayer.Y - z.Y
						if yDist*yDist > 0 { // Just checking for positive magnitude
							z.Y += yDist / int64(math.Abs(float64(yDist)))
						}
					}
				} else { // Lose track of the player if it is too far
					z.FollowingPlayerId = -1
				}
			} else { // otherwise randomly move
				randomDirection := util.RandomIntBetween(0, 4)
				if randomDirection == 0 {
					z.Y--
				} else if randomDirection == 1 {
					z.Y++
				} else if randomDirection == 2 {
					z.X--
				} else if randomDirection == 3 {
					z.X++
				}
			}
			// try to attack a nearby player
			for _, otherPlayer := range gs.Players {

				xDist := otherPlayer.X - z.X
				yDist := otherPlayer.Y - z.Y

				// Check for adjacency
				if xDist*xDist <= 1 && yDist*yDist <= 1 {
					otherPlayer.HP -= z.Damage
					if otherPlayer.HP <= 0 {
						gs.Logger.Printf("%s was struck down", otherPlayer.Username)

						// TODO - Drop stuff and do a respawn
						otherPlayer.X = 0
						otherPlayer.Y = 0
						otherPlayer.Gold = 0
						otherPlayer.HP = 10

						for _, player := range gs.Players {
							player.Conn.Write(packets.BuildUpdateSelfPlayerPacket(otherPlayer.ToBytes()))
						}
					} else {
						gs.Logger.Printf("%s was struck, %d HP remains", otherPlayer.Username, otherPlayer.HP)
						// send update to all players
						for _, player := range gs.Players {
							player.Conn.Write(packets.BuildUpdateSelfPlayerPacket(otherPlayer.ToBytes()))
						}
					}
					break
				}
			}

			// send zombie updates to all players
			for _, otherPlayer := range gs.Players {
				otherPlayer.Conn.Write(packets.BuildUpdateZombiePacket(z.ToBytes()))
			}

		}
	}
	gs.Mutex.Unlock()
}

func (gs *GameState) UpdateCombat() {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	for _, player := range gs.Players {
		if player.Hostile {
			// Attack the first zombie you find in range
			for _, zombie := range gs.Zombies {
				xdist := zombie.X - player.X
				ydist := zombie.Y - player.Y

				// Must be on an adjacent or the same tile
				// Diagonal works too
				if xdist*xdist <= 1 && ydist*ydist <= 1 {
					zombie.HP -= player.ST
					if zombie.HP <= 0 {
						gs.Logger.Printf("Zombie was struck down")
						for _, player := range gs.Players {
							player.Conn.Write(packets.BuildRemoveZombiePacket(zombie.ID))
						}
						delete(gs.Zombies, zombie.ID)
					} else {
						gs.Logger.Printf("Zombie was struck, %d HP remains", zombie.HP)
						// send zombie update to all players
						for _, otherPlayer := range gs.Players {
							otherPlayer.Conn.Write(packets.BuildUpdateZombiePacket(zombie.ToBytes()))
						}
					}
					break
				}
			}
		}
	}
}

func (gs *GameState) RunGameLoop() {
	for {
		dt := time.Duration((1.0 / float64(gs.TPS)) * 1000)
		gs.UpdateCombat()
		gs.UpdateZombies()

		time.Sleep(dt * time.Millisecond)
	}
}
