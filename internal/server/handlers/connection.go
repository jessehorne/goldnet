package handlers

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/shared"
	"github.com/jessehorne/goldnet/internal/util"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"net"

	"github.com/jessehorne/goldnet/internal/game/inventory"

	"github.com/jessehorne/goldnet/internal/game"
)

func ServerUserJoinHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user joined with a ID of", playerID)

	// add player to gamestates list of players
	newPlayer := game.NewPlayer(playerID, 0, 0, nil, conn)
	gs.AddPlayer(newPlayer)

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
		util.Send(conn, p)
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
		util.Send(conn, p)
	})
	newPlayer.Inventory.AddItem(clueNote)

	// send zombies to player
	for _, z := range gs.Zombies {
		zPacket := &packets.UpdateZombie{
			Type:              shared.PacketUpdateZombie,
			Id:                z.ID,
			X:                 z.X,
			Y:                 z.Y,
			Hp:                z.HP,
			Damage:            z.Damage,
			GoldDrop:          z.GoldDropAmt,
			FollowingPlayerId: z.FollowingPlayerId,
		}
		zData, zerr := proto.Marshal(zPacket)
		if zerr != nil {
			gs.Logger.Println(zerr)
			continue
		}
		util.Send(newPlayer.Conn, zData)
	}

	// let every player know they joined
	others := []*game.Player{}
	for _, p := range gs.Players {
		if p == nil {
			continue
		}
		if p.ID == playerID {
			continue
		}
		others = append(others, p)

		pJoined := &packets.PlayerJoined{
			Type: shared.PacketPlayerJoined,
			Id:   newPlayer.ID,
			X:    newPlayer.X,
			Y:    newPlayer.Y,
		}
		pData, perr := proto.Marshal(pJoined)
		if perr != nil {
			gs.Logger.Println(perr)
			continue
		}
		util.Send(p.Conn, pData)
	}

	var otherPlayers []*packets.UpdatePlayer
	for _, player := range others {
		up := &packets.UpdatePlayer{
			Id:       player.ID,
			Username: player.Username,
			X:        player.X,
			Y:        player.Y,
			Hp:       player.HP,
			Hostile:  player.Hostile,
		}
		otherPlayers = append(otherPlayers, up)
	}

	// send self join packet to player with their ID
	selfUpdate := &packets.SelfJoin{
		Type: shared.PacketPlayerSelfJoined,
		Self: &packets.UpdatePlayer{
			Id:        newPlayer.ID,
			Username:  newPlayer.Username,
			X:         newPlayer.X,
			Y:         newPlayer.Y,
			Gold:      newPlayer.Gold,
			Hp:        newPlayer.HP,
			St:        newPlayer.ST,
			Hostile:   newPlayer.Hostile,
			Inventory: newPlayer.Inventory.ToBytes(),
		},
		Others: otherPlayers,
	}
	selfUpdateData, selfUpdateError := proto.Marshal(selfUpdate)
	if selfUpdateError != nil {
		gs.Logger.Println(selfUpdateError)
		return
	}
	util.Send(conn, selfUpdateData)

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
	util.Send(conn, chunksPacketData)
}

func ServerUserDisconnectedHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user disconnected")

	// remove player from gamestate
	gs.RemovePlayer(playerID)

	// let everyone know they left
	dp := &packets.PlayerDisconnected{
		Type: shared.PacketPlayerDisconnected,
		Id:   playerID,
	}

	dpData, dpErr := proto.Marshal(dp)
	if dpErr != nil {
		gs.Logger.Println(dpErr)
		return
	}
	for _, p := range gs.Players {
		if p == nil {
			continue
		}
		util.Send(p.Conn, dpData)
	}
}
