package game

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/game/components"
	"github.com/jessehorne/goldnet/internal/shared"

	"github.com/jessehorne/goldnet/internal/util"
)

type GameState struct {
	EntityCounter      int64
	PlayerCounter      int64
	Players            map[components.EntityId]*components.PlayerComponent
	Zombies            map[components.EntityId]*Zombie
	SpriteComponents   map[components.EntityId]*components.SpriteComponent
	PositionComponents map[components.EntityId]*components.PositionComponent
	Mutex              sync.Mutex
	Logger             *log.Logger
	Chunks             map[int64]map[int64]*Chunk
	IntStore           map[string]int64
	TPS                int // ticks per second
}

func NewGameState() *GameState {
	return &GameState{
		Logger:             log.New(os.Stdout, "[GoldNet] (GameState) ", log.Ldate|log.Ltime),
		Players:            map[components.EntityId]*components.PlayerComponent{},
		Zombies:            map[components.EntityId]*Zombie{},
		SpriteComponents:   map[components.EntityId]*components.SpriteComponent{},
		PositionComponents: map[components.EntityId]*components.PositionComponent{},
		EntityCounter:      0,
		PlayerCounter:      0,
		Chunks:             map[int64]map[int64]*Chunk{},
		IntStore:           map[string]int64{},
		TPS:                10, // ticks per second
	}
}

func (gs *GameState) NextEntityId() components.EntityId {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	gs.EntityCounter += 1
	return components.EntityId(gs.EntityCounter)
}

func (gs *GameState) GetPlayer(playerID int64) *components.PlayerComponent {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	p, ok := gs.Players[components.EntityId(playerID)]
	if !ok {
		return nil
	}
	return p
}

func (gs *GameState) NextPlayerID() components.EntityId {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	gs.PlayerCounter += 1
	return components.EntityId(gs.PlayerCounter)
}

func (gs *GameState) AddPlayer(p *components.PlayerComponent) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	gs.Players[p.ID] = p
	gs.SpriteComponents[p.ID] = components.NewSpriteComponent('@', tcell.ColorWhite, tcell.ColorBlack)
	gs.PositionComponents[p.ID] = components.NewPositionComponent(0, 0)
}

func (gs *GameState) RemovePlayer(playerID components.EntityId) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	delete(gs.Players, playerID)
}

func (gs *GameState) UpdatePlayerChunks(playerID components.EntityId, x, y int64) {
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

func (gs *GameState) MovePlayer(playerID components.EntityId, x, y int64) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	p := gs.Players[playerID]
	position := gs.PositionComponents[playerID]
	if p != nil {
		position.X = x
		position.Y = y
		p.OldChunkX = position.X / CHUNK_W
		p.OldChunkY = position.Y / CHUNK_H
	}
}

func (gs *GameState) HandlePlayerAction(player *components.PlayerComponent, action int32) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	if player != nil {
		position := gs.PositionComponents[player.ID]
		if action == shared.ActionMoveLeft {
			position.X--
		} else if action == shared.ActionMoveRight {
			position.X++
		} else if action == shared.ActionMoveUp {
			position.Y--
		} else if action == shared.ActionMoveDown {
			position.Y++
		}
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

func (gs *GameState) GetTerrainAtCoords(x, y int64) byte {
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
	return shared.GetTerrainType(c.Stack[modY][modX][0])
}

func (gs *GameState) GetChunksAroundPlayer(p *components.PlayerComponent) ([]*Chunk, []*Chunk) {
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

func (gs *GameState) GetPlayersAroundPlayer(p *components.PlayerComponent) []*components.PlayerComponent {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	var players []*components.PlayerComponent
	if p == nil {
		return players
	}

	position := gs.PositionComponents[p.ID]
	if position == nil {
		// TODO - Maybe shouldn't fail silently
		return players
	}

	for _, otherPlayer := range gs.Players {
		if otherPlayer == nil {
			continue
		}
		if p.ID == otherPlayer.ID {
			continue
		}
		otherPosition := gs.PositionComponents[otherPlayer.ID]

		dis := util.Distance(position.X, position.Y, otherPosition.X, otherPosition.Y)
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

	// // update zombie
	// gs.Mutex.Lock()
	// for _, z := range gs.Zombies {
	// 	// handle movement
	// 	doesMove := util.RandomIntBetween(0, 10) < 2
	// 	if doesMove {
	// 		// If not currently following a player, pick one
	// 		if z.FollowingPlayerId == -1 {
	// 			for _, player := range gs.Players {
	// 				playerPosition := gs.PositionComponents[player.ID]
	// 				if util.Distance(z.X, z.Y, playerPosition.X, playerPosition.Y) < ZOMBIE_FOLLOW_RANGE {
	// 					z.FollowingPlayerId = player.ID
	// 				}
	// 			}
	// 		}

	// 		// If we are now following a player, move towards it
	// 		if z.FollowingPlayerId != -1 {
	// 			if z == nil {
	// 				continue
	// 			}
	// 			followingPosition := gs.PositionComponents[z.FollowingPlayerId]
	// 			if followingPosition == nil {
	// 				continue
	// 			}

	// 			// Follow player if close enough
	// 			if util.Distance(z.X, z.Y, followingPosition.X, followingPosition.Y) < ZOMBIE_FOLLOW_RANGE {
	// 				direction := util.RandomIntBetween(0, 2)
	// 				if direction == 0 {
	// 					xDist := followingPosition.X - z.X
	// 					if xDist*xDist > 0 { // Just checking for positive magnitude
	// 						z.X += xDist / int64(math.Abs(float64(xDist)))
	// 					}
	// 				} else {
	// 					yDist := followingPosition.Y - z.Y
	// 					if yDist*yDist > 0 { // Just checking for positive magnitude
	// 						z.Y += yDist / int64(math.Abs(float64(yDist)))
	// 					}
	// 				}
	// 			} else { // Lose track of the player if it is too far
	// 				z.FollowingPlayerId = -1
	// 			}
	// 		} else { // otherwise randomly move
	// 			randomDirection := util.RandomIntBetween(0, 4)
	// 			if randomDirection == 0 {
	// 				z.Y--
	// 			} else if randomDirection == 1 {
	// 				z.Y++
	// 			} else if randomDirection == 2 {
	// 				z.X--
	// 			} else if randomDirection == 3 {
	// 				z.X++
	// 			}
	// 		}
	// 		// try to attack a nearby player
	// 		timePerAttack := 1500.0
	// 		canAttackAt := z.LastAttackTime.Add(time.Duration(timePerAttack) * time.Millisecond)
	// 		if canAttackAt.Before(time.Now()) {
	// 			for _, otherPlayer := range gs.Players {
	// 				otherPlayerPosition := gs.PositionComponents[otherPlayer.ID]
	// 				xDist := otherPlayerPosition.X - z.X
	// 				yDist := otherPlayerPosition.Y - z.Y

	// 				// Check for adjacency
	// 				if xDist*xDist <= 1 && yDist*yDist <= 1 {
	// 					z.LastAttackTime = time.Now()
	// 					otherPlayer.HP -= z.Damage

	// 					msg := fmt.Sprintf("(GAME) You were struck by zombie for %d HP", z.Damage)
	// 					msgPacket := &packets.Message{
	// 						Type: shared.PacketSendMessage,
	// 						Data: msg,
	// 					}
	// 					p, perr := proto.Marshal(msgPacket)
	// 					if perr != nil {
	// 						gs.Logger.Println(perr)
	// 						return
	// 					}
	// 					util.Send(otherPlayer.Conn, p)

	// 					if otherPlayer.HP <= 0 {
	// 						msgPacket = &packets.Message{
	// 							Type: shared.PacketSendMessage,
	// 							Data: "(GAME) YOU WERE STRUCK DOWN BY ZOMBIE",
	// 						}
	// 						p, perr = proto.Marshal(msgPacket)
	// 						if perr != nil {
	// 							gs.Logger.Println(perr)
	// 							return
	// 						}
	// 						util.Send(otherPlayer.Conn, p)

	// 						// TODO - Drop stuff and do a respawn
	// 						otherPlayerPosition.X = 0
	// 						otherPlayerPosition.Y = 0
	// 						otherPlayer.Gold = 0
	// 						otherPlayer.HP = 10
	// 					}

	// 					// send update to all players
	// 					upp := &packets.UpdatePlayer{
	// 						Type:      shared.PacketUpdatePlayer,
	// 						Id:        int64(otherPlayer.ID),
	// 						Username:  otherPlayer.Username,
	// 						Gold:      otherPlayer.Gold,
	// 						Hp:        otherPlayer.HP,
	// 						St:        otherPlayer.ST,
	// 						Hostile:   otherPlayer.Hostile,
	// 						Inventory: otherPlayer.Inventory.ToBytes(),
	// 					}
	// 					uppData, uppDataErr := proto.Marshal(upp)
	// 					if uppDataErr != nil {
	// 						gs.Logger.Println(uppDataErr)
	// 						return
	// 					}
	// 					for _, player := range gs.Players {
	// 						util.Send(player.Conn, uppData)
	// 					}

	// 					break
	// 				}
	// 			}
	// 		}

	// 		// send zombie updates to all players
	// 		zPacket := &packets.UpdateZombie{
	// 			Type:              shared.PacketUpdateZombie,
	// 			Id:                int64(z.ID),
	// 			X:                 z.X,
	// 			Y:                 z.Y,
	// 			Hp:                z.HP,
	// 			Damage:            z.Damage,
	// 			GoldDrop:          z.GoldDropAmt,
	// 			FollowingPlayerId: int64(z.FollowingPlayerId),
	// 		}
	// 		zData, zerr := proto.Marshal(zPacket)
	// 		if zerr != nil {
	// 			gs.Logger.Println(zerr)
	// 			continue
	// 		}
	// 		for _, otherPlayer := range gs.Players {
	// 			util.Send(otherPlayer.Conn, zData)
	// 		}

	// 	}
	// }
	// gs.Mutex.Unlock()
}

func (gs *GameState) UpdateCombat() {
	// gs.Mutex.Lock()
	// defer gs.Mutex.Unlock()

	// for _, player := range gs.Players {
	// 	timePerAttack := 1000.0 / player.AttackSpeed
	// 	canAttackAt := player.LastAttackTime.Add(time.Duration(timePerAttack) * time.Millisecond)
	// 	if player.Hostile && canAttackAt.Before(time.Now()) {
	// 		// Attack the first zombie you find in range
	// 		for _, zombie := range gs.Zombies {
	// 			xdist := zombie.X - player.X
	// 			ydist := zombie.Y - player.Y

	// 			// Must be on an adjacent or the same tile
	// 			// Diagonal works too
	// 			if xdist*xdist <= 1 && ydist*ydist <= 1 {
	// 				player.LastAttackTime = time.Now()
	// 				zombie.HP -= player.ST
	// 				if zombie.HP <= 0 {
	// 					msgPacket := &packets.Message{
	// 						Type: shared.PacketSendMessage,
	// 						Data: "(GAME) You struck the zombie down",
	// 					}
	// 					p, perr := proto.Marshal(msgPacket)
	// 					if perr != nil {
	// 						gs.Logger.Println(perr)
	// 						return
	// 					}
	// 					util.Send(player.Conn, p)

	// 					rz := &packets.RemoveZombie{
	// 						Type: shared.PacketRemoveZombie,
	// 						Id:   zombie.ID,
	// 					}
	// 					rzData, rzErr := proto.Marshal(rz)
	// 					if rzErr != nil {
	// 						gs.Logger.Println(rzErr)
	// 						return
	// 					}

	// 					for _, otherPlayer := range gs.Players {
	// 						util.Send(otherPlayer.Conn, rzData)
	// 					}
	// 					delete(gs.Zombies, zombie.ID)
	// 				} else {
	// 					msg := fmt.Sprintf("You struck the zombie for %d HP", player.ST)
	// 					msgPacket := &packets.Message{
	// 						Type: shared.PacketSendMessage,
	// 						Data: msg,
	// 					}
	// 					p, perr := proto.Marshal(msgPacket)
	// 					if perr != nil {
	// 						gs.Logger.Println(perr)
	// 						return
	// 					}
	// 					util.Send(player.Conn, p)

	// 					// send zombie update to all players
	// 					zPacket := &packets.UpdateZombie{
	// 						Type:              shared.PacketUpdateZombie,
	// 						Id:                zombie.ID,
	// 						X:                 zombie.X,
	// 						Y:                 zombie.Y,
	// 						Hp:                zombie.HP,
	// 						Damage:            zombie.Damage,
	// 						GoldDrop:          zombie.GoldDropAmt,
	// 						FollowingPlayerId: zombie.FollowingPlayerId,
	// 					}
	// 					zData, zerr := proto.Marshal(zPacket)
	// 					if zerr != nil {
	// 						gs.Logger.Println(zerr)
	// 						continue
	// 					}
	// 					for _, otherPlayer := range gs.Players {
	// 						util.Send(otherPlayer.Conn, zData)
	// 					}
	// 				}

	// 				goto endattackattempt
	// 			}
	// 		}

	// 		for _, otherPlayer := range gs.Players {

	// 			// Suicide watch
	// 			if otherPlayer.ID == player.ID {
	// 				continue
	// 			}

	// 			xdist := otherPlayer.X - player.X
	// 			ydist := otherPlayer.Y - player.Y

	// 			// Must be on an adjacent or the same tile
	// 			// Diagonal works too
	// 			if xdist*xdist <= 1 && ydist*ydist <= 1 {
	// 				player.LastAttackTime = time.Now()
	// 				otherPlayer.HP -= player.ST

	// 				msg2 := fmt.Sprintf("You struck %s for %d HP", otherPlayer.Username, player.ST)
	// 				msgPacket := &packets.Message{
	// 					Type: shared.PacketSendMessage,
	// 					Data: msg2,
	// 				}
	// 				p, perr := proto.Marshal(msgPacket)
	// 				if perr != nil {
	// 					gs.Logger.Println(perr)
	// 					return
	// 				}
	// 				util.Send(player.Conn, p)

	// 				msg := fmt.Sprintf("You were struck by %s for %d HP", player.Username, player.ST)
	// 				msgPacket = &packets.Message{
	// 					Type: shared.PacketSendMessage,
	// 					Data: msg,
	// 				}
	// 				p, perr = proto.Marshal(msgPacket)
	// 				if perr != nil {
	// 					gs.Logger.Println(perr)
	// 					return
	// 				}
	// 				util.Send(otherPlayer.Conn, p)

	// 				if otherPlayer.HP <= 0 {
	// 					msg = fmt.Sprintf("You struck down %s", otherPlayer.Username)
	// 					msgPacket = &packets.Message{
	// 						Type: shared.PacketSendMessage,
	// 						Data: msg,
	// 					}
	// 					p, perr = proto.Marshal(msgPacket)
	// 					if perr != nil {
	// 						gs.Logger.Println(perr)
	// 						return
	// 					}
	// 					util.Send(player.Conn, p)

	// 					msg2 = fmt.Sprintf("YOU WERE STRUCK DOWN BY %s", player.Username)
	// 					msgPacket = &packets.Message{
	// 						Type: shared.PacketSendMessage,
	// 						Data: msg2,
	// 					}
	// 					p, perr = proto.Marshal(msgPacket)
	// 					if perr != nil {
	// 						gs.Logger.Println(perr)
	// 						return
	// 					}
	// 					util.Send(otherPlayer.Conn, p)

	// 					// TODO - Drop stuff and do a respawn
	// 					otherPlayer.X = 0
	// 					otherPlayer.Y = 0
	// 					otherPlayer.Gold = 0
	// 					otherPlayer.HP = 10
	// 				}

	// 				// send update to all players
	// 				upp := &packets.UpdatePlayer{
	// 					Type:      shared.PacketUpdatePlayer,
	// 					Id:        otherPlayer.ID,
	// 					Username:  otherPlayer.Username,
	// 					X:         otherPlayer.X,
	// 					Y:         otherPlayer.Y,
	// 					Gold:      otherPlayer.Gold,
	// 					Hp:        otherPlayer.HP,
	// 					St:        otherPlayer.ST,
	// 					Hostile:   otherPlayer.Hostile,
	// 					Inventory: otherPlayer.Inventory.ToBytes(),
	// 				}
	// 				uppData, uppDataErr := proto.Marshal(upp)
	// 				if uppDataErr != nil {
	// 					gs.Logger.Println(uppDataErr)
	// 					return
	// 				}
	// 				for _, pl := range gs.Players {
	// 					util.Send(pl.Conn, uppData)
	// 				}

	// 				goto endattackattempt
	// 			}
	// 		}

	// 	endattackattempt:
	// 		continue
	// 	}
	// }
}

func (gs *GameState) UseItem(p *components.PlayerComponent, itemID int64) {
	for _, item := range p.Inventory.Items {
		if item.GetID() == itemID {
			item.TriggerUse()
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
