package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
	"net"
	"time"
)

func ServerActionHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	action := data[0]

	p := gs.GetPlayer(playerID)
	if p == nil {
		return
	}

	if packets.IsMovementAction(action) {
		mod := (1 / float64(p.Speed)) * 1000

		b := gs.GetBelowBlockAtCoords(p.X, p.Y)
		if shared.GetTerrainBelow(b) == shared.TerrainWater {
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

			gs.HandlePlayerAction(p, action)

			// send the updated position to the player
			conn.Write(packets.BuildMovePacket(p.ID, p.X, p.Y))

			// send movement to other nearby players
			nearbyPlayers := gs.GetPlayersAroundPlayer(p)
			for _, other := range nearbyPlayers {
				other.Conn.Write(packets.BuildMovePacket(p.ID, p.X, p.Y))
			}
		} else {
			conn.Write(packets.BuildMovePacket(p.ID, p.X, p.Y))
		}

		// send chunks if players chunk has updated
		newChunkX := p.X / game.CHUNK_W
		newChunkY := p.Y / game.CHUNK_H
		if newChunkX != p.OldChunkX || newChunkY != p.OldChunkY {
			p.OldChunkX = p.X / game.CHUNK_W
			p.OldChunkY = p.Y / game.CHUNK_H
			nearbyChunks, newlyGenerated := gs.GetChunksAroundPlayer(p)

			// zombie spawning
			for _, c := range newlyGenerated {
				shouldCreateZombie := util.RandomIntBetween(0, 16) == 0
				if shouldCreateZombie {
					newZombie := game.NewZombie(c.X*game.CHUNK_W, c.Y*game.CHUNK_H)
					gs.Logger.Println("added zombie with ID", newZombie.ID)
					gs.Zombies[newZombie.ID] = newZombie

					// send new zombie to all players
					zombieBytes := newZombie.ToBytes()
					newZombiePacket := packets.BuildNewZombiePacket(zombieBytes)
					for _, otherPlayer := range gs.Players {
						otherPlayer.Conn.Write(newZombiePacket)
					}
				}
			}

			var chunkData []byte
			for _, c := range nearbyChunks {
				chunkData = append(chunkData, c.ToBytes()...)
			}
			chunkPacket := packets.BuildChunksPacket(int64(len(nearbyChunks)), chunkData)
			conn.Write(chunkPacket)
		}
	}
}
