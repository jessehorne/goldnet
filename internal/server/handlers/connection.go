package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/game/inventory"

	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

func ServerUserJoinHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user joined with ID of", playerID)

	// add player to gamestates list of players
	newPlayer := game.NewPlayer(playerID, 0, 0, nil, conn)
	gs.AddPlayer(newPlayer)

	// add a welcome note to the players inventory
	welcomeNote := inventory.NewNote(1, "a clean envelope", "welcome!")
	welcomeNote.SetUseCallback(func() {
		conn.Write(packets.BuildMessagePacket(-1, "The note says 'Welcome!'"))
	})
	newPlayer.Inventory.AddItem(welcomeNote)

	clueNote := inventory.NewNote(1, "a dirty envelope", "here's a clue...")
	clueNote.SetUseCallback(func() {
		conn.Write(packets.BuildMessagePacket(-1, "The note says 'here's a clue...'"))
	})
	newPlayer.Inventory.AddItem(clueNote)

	// send zombies to player
	for _, z := range gs.Zombies {
		newPlayer.Conn.Write(packets.BuildUpdateZombiePacket(z.ToBytes()))
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
		p.Conn.Write(packets.BuildPlayerJoinedPacket(newPlayer.ID, newPlayer.X, newPlayer.Y))
	}

	othersData := util.Int64ToBytes(int64(len(others)))

	// add other players and their positions
	for _, player := range others {
		othersData = append(othersData, util.Int64ToBytes(player.ID)...)
		othersData = append(othersData, util.Int64ToBytes(player.X)...)
		othersData = append(othersData, util.Int64ToBytes(player.Y)...)
	}

	// send self join packet to player with their ID
	conn.Write(packets.BuildPlayerSelfJoinedPacket(playerID, 0, 0, othersData, newPlayer.Inventory.ToBytes()))

	// send nearby chunks to player
	nearbyChunks, _ := gs.GetChunksAroundPlayer(newPlayer)
	var chunkData []byte
	for _, c := range nearbyChunks {
		chunkData = append(chunkData, c.ToBytes()...)
	}
	conn.Write(packets.BuildChunksPacket(int64(len(nearbyChunks)), chunkData))
}

func ServerUserDisconnectedHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	gs.Logger.Println("[PACKET] user disconnected")

	// remove player from gamestate
	gs.RemovePlayer(playerID)

	// let everyone know they left
	for _, p := range gs.Players {
		if p == nil {
			continue
		}
		p.Conn.Write(packets.BuildPlayerDisconnectedPacket(playerID))
	}
}
