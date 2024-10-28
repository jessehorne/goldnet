package handlers

import (
	"net"
	"time"

	packets "github.com/jessehorne/goldnet/packets/dist"
	packetscomponents "github.com/jessehorne/goldnet/packets/dist/components"
	"google.golang.org/protobuf/proto"

	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/game/components"
	"github.com/jessehorne/goldnet/internal/shared"
	"github.com/jessehorne/goldnet/internal/util"
)

func ServerActionHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	var action packets.Action
	err := proto.Unmarshal(data, &action)
	if err != nil {
		gs.Logger.Println(err)
		return
	}

	p := gs.GetPlayer(playerID)
	if p == nil {
		return
	}

	position := gs.PositionComponents[p.ID]

	if shared.IsMovementAction(action.Action) {
		mod := (1 / float64(p.Speed)) * 1000

		b := gs.GetTerrainAtCoords(position.X, position.Y)
		if b == shared.TerrainWater {
			mod = mod * 4
		}

		canMoveAt := p.LastMovementTime.Add(time.Duration(mod) * time.Millisecond)
		canMove := true
		if time.Now().Before(canMoveAt) {
			canMove = false
			return
		}

		if canMove {
			p.LastMovementTime = time.Now()

			gs.HandlePlayerAction(p, action.Action)

			upp := &packetscomponents.UpdatePosition{
				Type:     shared.PacketUpdatePosition,
				EntityId: int64(p.ID),
				X:        position.X,
				Y:        position.Y,
			}
			uppData, uppDataErr := proto.Marshal(upp)
			if uppDataErr != nil {
				gs.Logger.Println(uppDataErr)
				return
			}
			util.Send(conn, uppData)

			// send movement to other nearby players
			nearbyPlayers := gs.GetPlayersAroundPlayer(p)
			for _, other := range nearbyPlayers {
				util.Send(other.Conn, uppData)
			}
		} else {
			// send the updated position to the player
			upp := &packets.UpdatePlayer{
				Type:      shared.PacketUpdatePosition,
				Id:        int64(p.ID),
				Username:  p.Username,
				Gold:      p.Gold,
				Hp:        p.HP,
				St:        p.ST,
				Hostile:   p.Hostile,
				Inventory: p.Inventory.ToBytes(),
			}
			uppData, uppDataErr := proto.Marshal(upp)
			if uppDataErr != nil {
				gs.Logger.Println(uppDataErr)
				return
			}
			util.Send(conn, uppData)
		}

		// send chunks if players chunk has updated
		newChunkX := position.X / game.CHUNK_W
		newChunkY := position.Y / game.CHUNK_H
		if newChunkX != p.OldChunkX || newChunkY != p.OldChunkY {
			p.OldChunkX = position.X / game.CHUNK_W
			p.OldChunkY = position.Y / game.CHUNK_H
			nearbyChunks, newlyGenerated := gs.GetChunksAroundPlayer(p)

			// zombie spawning
			for _, c := range newlyGenerated {
				shouldCreateZombie := util.RandomIntBetween(0, 16) == 0
				if shouldCreateZombie {
					newX := c.X * game.CHUNK_W
					newY := c.Y * game.CHUNK_H
					newZombie := components.NewZombieComponent(gs.NextEntityId(), newX, newY)
					gs.Logger.Println("added zombie with ID", newZombie.ID)
					gs.InitNewZombie(newZombie, newX, newY)
				}
			}

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

			util.Send(conn, chunksPacketData)
		}
	}
}

func ServerSetHostileHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	var sh packets.SetHostile
	err := proto.Unmarshal(data, &sh)
	if err != nil {
		gs.Logger.Println(err)
		return
	}

	p := gs.GetPlayer(playerID)
	if p == nil {
		return
	}

	p.Hostile = sh.Hostile

	gs.Logger.Println("Toggled Hostile")

	for _, player := range gs.PlayerComponents {
		setHostile := &packets.SetHostile{
			Type:     shared.PacketSetHostile,
			PlayerID: int64(p.ID),
			Hostile:  p.Hostile,
		}
		setHostileData, setHostileDataErr := proto.Marshal(setHostile)
		if setHostileDataErr != nil {
			gs.Logger.Println(err)
			continue
		}
		util.Send(player.Conn, setHostileData)
	}
}
