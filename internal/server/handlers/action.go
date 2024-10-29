package handlers

import (
	"net"
	"time"

	"github.com/gdamore/tcell/v2"
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
			game.SendOneToOne(conn, gs, &packets.UpdatePlayer{
				Type:      shared.PacketUpdatePosition,
				Id:        int64(p.ID),
				Username:  p.Username,
				Gold:      p.Gold,
				Hp:        p.HP,
				St:        p.ST,
				Hostile:   p.Hostile,
				Inventory: p.Inventory.ToBytes(),
			})
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

			game.SendOneToOne(conn, gs, &packets.Chunks{
				Type: shared.PacketChunks,
				Data: chunkData,
			})
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

	backgroundColor := tcell.ColorBlack
	if sh.Hostile {
		backgroundColor = tcell.ColorRed
	}

	// Update hostile sprites
	game.SendOneToAll(gs, &packetscomponents.UpdateSprite{
		Type:       shared.PacketUpdateSprite,
		EntityId:   int64(p.ID),
		Character:  '@',
		Foreground: int64(tcell.ColorWhite),
		Background: int64(backgroundColor),
	})
}
