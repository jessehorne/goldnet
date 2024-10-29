package game

import (
	"fmt"
	"time"

	"github.com/jessehorne/goldnet/internal/shared"
	packets "github.com/jessehorne/goldnet/packets/dist"
	packetscomponents "github.com/jessehorne/goldnet/packets/dist/components"
)

func UpdateCombatSystem(gs *GameState) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	for _, player := range gs.PlayerComponents {
		timePerAttack := 1000.0 / player.AttackSpeed
		canAttackAt := player.LastAttackTime.Add(time.Duration(timePerAttack) * time.Millisecond)
		if player.Hostile && canAttackAt.Before(time.Now()) {
			playerPosition := gs.PositionComponents[player.ID]
			if playerPosition == nil {
				gs.Logger.Println("Player is missing a position component")
				continue
			}
			// Attack the first zombie you find in range
			for _, zombie := range gs.ZombieComponents {
				zombiePosition := gs.PositionComponents[zombie.ID]
				if zombiePosition == nil {
					gs.Logger.Println("Zombie is missing a position component")
					continue
				}
				xdist := zombiePosition.X - playerPosition.X
				ydist := zombiePosition.Y - playerPosition.Y

				// Must be on an adjacent, diagonal, or the same tile
				if xdist*xdist <= 1 && ydist*ydist <= 1 {
					player.LastAttackTime = time.Now()
					zombie.HP -= player.ST
					if zombie.HP <= 0 {
						SendOneToOne(player.Conn, gs, &packets.Message{
							Type: shared.PacketSendMessage,
							Data: "(GAME) You struck the zombie down",
						})
						SendOneToAll(gs, &packets.RemoveZombie{
							Type: shared.PacketRemoveZombie,
							Id:   int64(zombie.ID),
						})
						gs.RemoveZombie(zombie.ID)
					} else {
						msg := fmt.Sprintf("You struck the zombie for %d HP", player.ST)
						SendOneToOne(player.Conn, gs, &packets.Message{
							Type: shared.PacketSendMessage,
							Data: msg,
						})

						// Send zombie update to all players
						SendOneToAll(gs, &packets.UpdateZombie{
							Type:              shared.PacketUpdateZombie,
							Id:                int64(zombie.ID),
							Hp:                zombie.HP,
							Damage:            zombie.Damage,
							GoldDrop:          zombie.GoldDropAmt,
							FollowingPlayerId: int64(zombie.FollowingPlayerId),
						})
					}

					goto endattackattempt
				}
			}

			for _, otherPlayer := range gs.PlayerComponents {

				// Suicide watch
				if otherPlayer.ID == player.ID {
					continue
				}

				otherPlayerPosition := gs.PositionComponents[otherPlayer.ID]
				if otherPlayerPosition == nil {
					gs.Logger.Println("Player is missing a position component")
					continue
				}

				xdist := otherPlayerPosition.X - playerPosition.X
				ydist := otherPlayerPosition.Y - playerPosition.Y

				// Must be on an adjacent or the same tile
				// Diagonal works too
				if xdist*xdist <= 1 && ydist*ydist <= 1 {
					player.LastAttackTime = time.Now()
					otherPlayer.HP -= player.ST

					msg1 := fmt.Sprintf("You struck %s for %d HP", otherPlayer.Username, player.ST)
					SendOneToOne(player.Conn, gs, &packets.Message{
						Type: shared.PacketSendMessage,
						Data: msg1,
					})

					msg2 := fmt.Sprintf("You were struck by %s for %d HP", player.Username, player.ST)
					SendOneToOne(otherPlayer.Conn, gs, &packets.Message{
						Type: shared.PacketSendMessage,
						Data: msg2,
					})

					if otherPlayer.HP <= 0 {
						msg := fmt.Sprintf("You struck down %s", otherPlayer.Username)
						SendOneToOne(player.Conn, gs, &packets.Message{
							Type: shared.PacketSendMessage,
							Data: msg,
						})

						msg2 = fmt.Sprintf("YOU WERE STRUCK DOWN BY %s", player.Username)
						SendOneToOne(otherPlayer.Conn, gs, &packets.Message{
							Type: shared.PacketSendMessage,
							Data: msg2,
						})

						// TODO - Drop stuff and do a respawn
						otherPlayer.Gold = 0
						otherPlayer.HP = 10

						SendOneToAll(gs, &packetscomponents.UpdatePosition{
							Type:     shared.PacketUpdatePosition,
							EntityId: int64(otherPlayer.ID),
							X:        otherPlayerPosition.X,
							Y:        otherPlayerPosition.Y,
						})
					}

					// send update to all players
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

					goto endattackattempt
				}
			}

		endattackattempt:
			continue
		}
	}
}
