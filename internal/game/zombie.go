package game

import (
	"fmt"
	"math"
	"time"

	"github.com/jessehorne/goldnet/internal/game/components"
	"github.com/jessehorne/goldnet/internal/shared"
	"github.com/jessehorne/goldnet/internal/util"
	packets "github.com/jessehorne/goldnet/packets/dist"
	packetscomponents "github.com/jessehorne/goldnet/packets/dist/components"
)

func UpdateZombiesSystem(gs *GameState) {
	// update zombie
	gs.Mutex.Lock()
	for zombieId, z := range gs.ZombieComponents {

		zombiePosition := gs.PositionComponents[zombieId]
		if zombiePosition == nil {
			gs.Logger.Printf("Missing position component for zombie")
			return
		}

		// handle movement
		doesMove := util.RandomIntBetween(0, 10) < 2
		if doesMove {
			// If not currently following a player, pick one
			if z.FollowingPlayerId == -1 {
				for _, player := range gs.PlayerComponents {
					playerPosition := gs.PositionComponents[player.ID]
					if util.Distance(zombiePosition.X, zombiePosition.Y, playerPosition.X, playerPosition.Y) <
						components.ZOMBIE_FOLLOW_RANGE {
						z.FollowingPlayerId = player.ID
					}
				}
			}

			// If we are now following a player, move towards it
			if z.FollowingPlayerId != -1 {
				if z == nil {
					continue
				}
				followingPosition := gs.PositionComponents[z.FollowingPlayerId]
				if followingPosition == nil {
					continue
				}

				// Follow player if close enough
				if util.Distance(zombiePosition.X, zombiePosition.Y, followingPosition.X, followingPosition.Y) <
					components.ZOMBIE_FOLLOW_RANGE {
					direction := util.RandomIntBetween(0, 2)
					if direction == 0 {
						xDist := followingPosition.X - zombiePosition.X
						if xDist*xDist > 0 { // Just checking for positive magnitude
							zombiePosition.X += xDist / int64(math.Abs(float64(xDist)))
						}
					} else {
						yDist := followingPosition.Y - zombiePosition.Y
						if yDist*yDist > 0 { // Just checking for positive magnitude
							zombiePosition.Y += yDist / int64(math.Abs(float64(yDist)))
						}
					}
				} else { // Lose track of the player if it is too far
					z.FollowingPlayerId = -1
				}
			} else { // otherwise randomly move
				randomDirection := util.RandomIntBetween(0, 4)
				if randomDirection == 0 {
					zombiePosition.Y--
				} else if randomDirection == 1 {
					zombiePosition.Y++
				} else if randomDirection == 2 {
					zombiePosition.X--
				} else if randomDirection == 3 {
					zombiePosition.X++
				}
			}
			// try to attack a nearby player
			timePerAttack := 1500.0 // TODO - This should probably be a stat
			canAttackAt := z.LastAttackTime.Add(time.Duration(timePerAttack) * time.Millisecond)
			if canAttackAt.Before(time.Now()) {
				for _, otherPlayer := range gs.PlayerComponents {
					otherPlayerPosition := gs.PositionComponents[otherPlayer.ID]
					xDist := otherPlayerPosition.X - zombiePosition.X
					yDist := otherPlayerPosition.Y - zombiePosition.Y

					// Check for adjacency
					if xDist*xDist <= 1 && yDist*yDist <= 1 {
						z.LastAttackTime = time.Now()
						otherPlayer.HP -= z.Damage

						msg := fmt.Sprintf("(GAME) You were struck by zombie for %d HP", z.Damage)
						SendOneToOne(otherPlayer.Conn, gs, &packets.Message{
							Type: shared.PacketSendMessage,
							Data: msg,
						})

						if otherPlayer.HP <= 0 {
							SendOneToOne(otherPlayer.Conn, gs, &packets.Message{
								Type: shared.PacketSendMessage,
								Data: "(GAME) YOU WERE STRUCK DOWN BY ZOMBIE",
							})

							// TODO - Drop stuff and do a respawn
							otherPlayerPosition.X = 0
							otherPlayerPosition.Y = 0
							otherPlayer.Gold = 0
							otherPlayer.HP = 10

							SendOneToAll(gs, &packetscomponents.UpdatePosition{
								Type:     shared.PacketUpdatePosition,
								EntityId: int64(otherPlayer.ID),
								X:        otherPlayerPosition.X,
								Y:        otherPlayerPosition.Y,
							})
						}

						// Send PlayerComponent update to all players
						SendOneToAll(gs, &packets.UpdatePlayer{
							Type:      shared.PacketUpdatePlayer,
							Id:        int64(otherPlayer.ID),
							Username:  otherPlayer.Username,
							Gold:      otherPlayer.Gold,
							Hp:        otherPlayer.HP,
							St:        otherPlayer.ST,
							Hostile:   otherPlayer.Hostile,
							Inventory: otherPlayer.Inventory.ToBytes(),
						})

						break
					}
				}
			}

			// send zombie component updates to all players
			SendOneToAll(gs, &packets.UpdateZombie{
				Type:              shared.PacketUpdateZombie,
				Id:                int64(z.ID),
				Hp:                z.HP,
				Damage:            z.Damage,
				GoldDrop:          z.GoldDropAmt,
				FollowingPlayerId: int64(z.FollowingPlayerId),
			})

			// Send zombie position updates to all players
			SendOneToAll(gs, &packetscomponents.UpdatePosition{
				Type:     shared.PacketUpdatePosition,
				EntityId: int64(z.ID),
				X:        zombiePosition.X,
				Y:        zombiePosition.Y,
			})
		}
	}
	gs.Mutex.Unlock()
}
