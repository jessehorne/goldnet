package game

func UpdateCombatSystem(gs *GameState) {
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
