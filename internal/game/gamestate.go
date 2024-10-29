package game

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/game/components"
	"github.com/jessehorne/goldnet/internal/game/inventory"
	"github.com/jessehorne/goldnet/internal/shared"
	packets "github.com/jessehorne/goldnet/packets/dist"
	packetscomponents "github.com/jessehorne/goldnet/packets/dist/components"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/jessehorne/goldnet/internal/util"
)

type GameState struct {
	EntityCounter      int64
	PlayerComponents   map[components.EntityId]*components.Player
	ZombieComponents   map[components.EntityId]*components.Zombie
	SpriteComponents   map[components.EntityId]*components.Sprite
	PositionComponents map[components.EntityId]*components.Position
	Mutex              sync.Mutex
	Logger             *log.Logger
	Chunks             map[int64]map[int64]*Chunk
	IntStore           map[string]int64
	TPS                int // ticks per second
}

func NewGameState() *GameState {
	return &GameState{
		Logger:             log.New(os.Stdout, "[GoldNet] (GameState) ", log.Ldate|log.Ltime),
		PlayerComponents:   map[components.EntityId]*components.Player{},
		ZombieComponents:   map[components.EntityId]*components.Zombie{},
		SpriteComponents:   map[components.EntityId]*components.Sprite{},
		PositionComponents: map[components.EntityId]*components.Position{},
		EntityCounter:      0,
		Chunks:             map[int64]map[int64]*Chunk{},
		IntStore:           map[string]int64{},
		TPS:                10, // ticks per second
	}
}

func (gs *GameState) NextEntityId() components.EntityId {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	gs.EntityCounter += 1

	gs.Logger.Println("Assigning entity ID: ", gs.EntityCounter)

	return components.EntityId(gs.EntityCounter)
}

func (gs *GameState) UpdatePositionComponent(entityId components.EntityId,
	position *components.Position) {
	gs.PositionComponents[entityId] = position

	positionUpdate := &packetscomponents.UpdatePosition{
		Type:     shared.PacketUpdatePosition,
		EntityId: int64(entityId),
		X:        position.X,
		Y:        position.Y,
	}
	existingUpdateData, existingUpdateError := proto.Marshal(positionUpdate)
	if existingUpdateError != nil {
		gs.Logger.Println(existingUpdateError)
		return
	}

	for _, player := range gs.PlayerComponents {
		util.Send(player.Conn, existingUpdateData)
	}
}

func (gs *GameState) UpdateSpriteComponent(entityId components.EntityId,
	sprite *components.Sprite) {
	gs.SpriteComponents[entityId] = sprite

	spriteUpdate := &packetscomponents.UpdateSprite{
		Type:       shared.PacketUpdateSprite,
		EntityId:   int64(entityId),
		Character:  sprite.Character,
		Foreground: int64(sprite.Foreground),
		Background: int64(sprite.Background),
	}
	existingUpdateData, existingUpdateError := proto.Marshal(spriteUpdate)
	if existingUpdateError != nil {
		gs.Logger.Println(existingUpdateError)
		return
	}

	for _, player := range gs.PlayerComponents {
		util.Send(player.Conn, existingUpdateData)
	}
}

func (gs *GameState) GetPlayer(playerID int64) *components.Player {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	p, ok := gs.PlayerComponents[components.EntityId(playerID)]
	if !ok {
		return nil
	}
	return p
}

func (gs *GameState) AddPlayer(p *components.Player) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	gs.PlayerComponents[p.ID] = p
}

func SendOneToOne[T protoreflect.ProtoMessage](conn net.Conn, gs *GameState, packet T) {
	data, err := proto.Marshal(packet)
	if err != nil {
		gs.Logger.Println(err)
		return
	}
	util.Send(conn, data)
}

func SendOneToAll[T protoreflect.ProtoMessage](gs *GameState, packet T) {
	data, err := proto.Marshal(packet)
	if err != nil {
		gs.Logger.Println(err)
		return
	}

	for _, player := range gs.PlayerComponents {
		util.Send(player.Conn, data)
	}
}

func SendAllToOne(conn net.Conn, gs *GameState) {
	// Learn of other players and teach them of myself
	for _, player := range gs.PlayerComponents {
		existingPlayerUpdate := &packets.UpdatePlayer{
			Type:     shared.PacketUpdatePlayer,
			Id:       int64(player.ID),
			Username: player.Username,
			Hp:       player.HP,
			Hostile:  player.Hostile,
		}
		existingUpdateData, existingUpdateError := proto.Marshal(existingPlayerUpdate)
		if existingUpdateError != nil {
			gs.Logger.Println(existingUpdateError)
			return
		}
		util.Send(conn, existingUpdateData)
	}

	// send zombies to player
	for _, z := range gs.ZombieComponents {
		zPacket := &packets.UpdateZombie{
			Type:              shared.PacketUpdateZombie,
			Id:                int64(z.ID),
			Hp:                z.HP,
			Damage:            z.Damage,
			GoldDrop:          z.GoldDropAmt,
			FollowingPlayerId: int64(z.FollowingPlayerId),
		}
		zData, zerr := proto.Marshal(zPacket)
		if zerr != nil {
			gs.Logger.Println(zerr)
			continue
		}
		util.Send(conn, zData)
	}

	// Send all position components
	for entityId, position := range gs.PositionComponents {
		existingPlayerUpdate := &packetscomponents.UpdatePosition{
			Type:     shared.PacketUpdatePosition,
			EntityId: int64(entityId),
			X:        position.X,
			Y:        position.Y,
		}
		SendOneToOne(conn, gs, existingPlayerUpdate)
	}

	// Send all sprite components
	for entityId, sprite := range gs.SpriteComponents {
		existingPlayerUpdate := &packetscomponents.UpdateSprite{
			Type:       shared.PacketUpdateSprite,
			EntityId:   int64(entityId),
			Character:  sprite.Character,
			Background: int64(sprite.Background),
			Foreground: int64(sprite.Foreground),
		}
		existingUpdateData, existingUpdateError := proto.Marshal(existingPlayerUpdate)
		if existingUpdateError != nil {
			gs.Logger.Println(existingUpdateError)
			return
		}
		util.Send(conn, existingUpdateData)
	}
}

func (gs *GameState) InitNewPlayer(newPlayer *components.Player) {
	gs.Mutex.Lock()
	gs.PlayerComponents[newPlayer.ID] = newPlayer
	gs.Mutex.Unlock()

	// add a welcome note to the players inventory
	welcomeNote := inventory.NewNote(1, "a clean envelope", "welcome!")
	welcomeNote.SetUseCallback(func() {
		msg := &packets.Message{
			Type: shared.PacketSendMessage,
			Data: fmt.Sprintf("%s %s", "(GAME)", "The note says 'Welcome!'"),
		}
		p, perr := proto.Marshal(msg)
		if perr != nil {
			gs.Logger.Println(perr)
			return
		}
		util.Send(newPlayer.Conn, p)
	})
	newPlayer.Inventory.AddItem(welcomeNote)

	clueNote := inventory.NewNote(1, "a dirty envelope", "here's a clue...")
	clueNote.SetUseCallback(func() {
		msg := &packets.Message{
			Type: shared.PacketSendMessage,
			Data: fmt.Sprintf("%s %s", "(GAME)", "The note says 'here's a clue...'"),
		}
		p, perr := proto.Marshal(msg)
		if perr != nil {
			gs.Logger.Println(perr)
			return
		}
		util.Send(newPlayer.Conn, p)
	})
	newPlayer.Inventory.AddItem(clueNote)

	pJoined := &packets.PlayerJoined{
		Type: shared.PacketPlayerJoined,
		Id:   int64(newPlayer.ID),
	}
	pData, perr := proto.Marshal(pJoined)
	if perr != nil {
		gs.Logger.Println(perr)
		return
	}

	// let every player know they joined
	for _, p := range gs.PlayerComponents {
		if p == nil {
			continue
		}
		if p.ID == newPlayer.ID {
			continue
		}
		util.Send(p.Conn, pData)
	}

	newPlayerUpdate := &packets.UpdatePlayer{
		Type:      shared.PacketUpdatePlayer,
		Id:        int64(newPlayer.ID),
		Username:  newPlayer.Username,
		Gold:      newPlayer.Gold,
		Hp:        newPlayer.HP,
		St:        newPlayer.ST,
		Hostile:   newPlayer.Hostile,
		Inventory: newPlayer.Inventory.ToBytes(),
	}

	// send self join packet to player with their ID
	selfJoin := &packets.SelfJoin{
		Type: shared.PacketPlayerSelfJoined,
		Self: newPlayerUpdate,
	}
	selfJoinData, selfJoinError := proto.Marshal(selfJoin)
	if selfJoinError != nil {
		gs.Logger.Println(selfJoinError)
		return
	}
	util.Send(newPlayer.Conn, selfJoinData)

	// Send full gamestate to the connecting player
	// TODO - Only send components for nearby chunks
	SendAllToOne(newPlayer.Conn, gs)

	// Send connecting player's component to everyone
	// TODO - Only send component to players in nearby chunks
	SendOneToAll(gs, newPlayerUpdate)

	// Add sprite component to the new player entity
	sprite := components.NewSpriteComponent('@', tcell.ColorWhite, tcell.ColorBlack)
	gs.UpdateSpriteComponent(newPlayer.ID, sprite)

	// Add position component to the new player entity
	position := components.NewPositionComponent(0, 0)
	gs.UpdatePositionComponent(newPlayer.ID, position)

	// send nearby chunks to player
	nearbyChunks, _ := gs.GetChunksAroundPlayer(newPlayer)
	chunkData := util.Int64ToBytes(int64(len(nearbyChunks)))
	for _, c := range nearbyChunks {
		chunkData = append(chunkData, c.ToBytes()...)
	}
	chunksPacket := &packets.Chunks{
		Type: shared.PacketChunks,
		Data: chunkData,
	}
	chunksPacketData, chunksPacketErr := proto.Marshal(chunksPacket)
	if chunksPacketErr != nil {
		gs.Logger.Println(chunksPacketErr)
		return
	}
	util.Send(newPlayer.Conn, chunksPacketData)
}

func (gs *GameState) RemovePlayer(playerID components.EntityId) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	// TODO - Maybe delete every component associated with the playerID entity
	// So we can't accidentally forget one
	delete(gs.PlayerComponents, playerID)
	delete(gs.SpriteComponents, playerID)
	delete(gs.PositionComponents, playerID)
}

func (gs *GameState) InitNewZombie(newZombie *components.Zombie, x, y int64) {
	gs.Mutex.Lock()
	gs.ZombieComponents[newZombie.ID] = newZombie
	gs.Mutex.Unlock()

	// Send new zombie component to all players
	SendOneToAll(gs, &packets.UpdateZombie{
		Type:              shared.PacketUpdateZombie,
		Id:                int64(newZombie.ID),
		Hp:                newZombie.HP,
		Damage:            newZombie.Damage,
		GoldDrop:          newZombie.GoldDropAmt,
		FollowingPlayerId: int64(newZombie.FollowingPlayerId),
	})

	// Add sprite component to the new zombie entity
	sprite := components.NewSpriteComponent('Z', tcell.ColorWhite, tcell.ColorBlack)
	gs.UpdateSpriteComponent(newZombie.ID, sprite)

	// Add position component to the new zombie entity
	position := components.NewPositionComponent(x, y)
	gs.UpdatePositionComponent(newZombie.ID, position)
}

func (gs *GameState) RemoveZombie(zombieID components.EntityId) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	// Maybe we should have a general `RemoveEntity` function
	// Current design allows mistakes
	delete(gs.ZombieComponents, zombieID)
	delete(gs.SpriteComponents, zombieID)
	delete(gs.PositionComponents, zombieID)
}

func (gs *GameState) UpdatePlayerChunks(playerID components.EntityId, x, y int64) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	_, ok := gs.PlayerComponents[playerID]
	if ok {
		gs.PlayerComponents[playerID].OldChunkX = x
		gs.PlayerComponents[playerID].OldChunkY = y
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
	p := gs.PlayerComponents[playerID]
	position := gs.PositionComponents[playerID]
	if p != nil {
		position.X = x
		position.Y = y
		p.OldChunkX = position.X / CHUNK_W
		p.OldChunkY = position.Y / CHUNK_H
	}
}

func (gs *GameState) HandlePlayerAction(player *components.Player, action int32) {
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

func (gs *GameState) GetChunksAroundPlayer(p *components.Player) ([]*Chunk, []*Chunk) {
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

func (gs *GameState) GetPlayersAroundPlayer(p *components.Player) []*components.Player {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	var players []*components.Player
	if p == nil {
		return players
	}

	position := gs.PositionComponents[p.ID]
	if position == nil {
		// TODO - Maybe shouldn't fail silently
		return players
	}

	for _, otherPlayer := range gs.PlayerComponents {
		if otherPlayer == nil {
			continue
		}
		if p.ID == otherPlayer.ID {
			continue
		}
		otherPosition := gs.PositionComponents[otherPlayer.ID]
		if otherPosition == nil {
			continue
		}

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

func (gs *GameState) UseItem(p *components.Player, itemID int64) {
	for _, item := range p.Inventory.Items {
		if item.GetID() == itemID {
			item.TriggerUse()
		}
	}
}

func (gs *GameState) RunGameLoop() {
	for {
		dt := time.Duration((1.0 / float64(gs.TPS)) * 1000)

		UpdateCombatSystem(gs)
		UpdateZombiesSystem(gs)

		// TODO: Time the update logic and pass the remaining time here
		time.Sleep(dt * time.Millisecond)
	}
}
