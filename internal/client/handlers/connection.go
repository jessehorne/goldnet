package handlers

import (
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/util"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"net"
)

func ClientPlayerSelfJoinedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var sj packets.SelfJoin
	err := proto.Unmarshal(data, &sj)
	if err != nil {
		gs.Logger.Println(err)
		return
	}

	p := game.NewPlayer(sj.Self.Id, sj.Self.X, sj.Self.Y, sj.Self.Inventory, nil)
	gs.AddPlayer(p)
	gs.SetIntStore("playerID", sj.Self.Id)
	g.Sidebar.UpdatePlayerStats(p)

	for _, op := range sj.Others {
		otherPlayer := &game.Player{
			ID:       op.Id,
			X:        op.X,
			Y:        op.Y,
			Username: op.Username,
			HP:       op.Hp,
			ST:       op.St,
			Hostile:  op.Hostile,
		}
		gs.AddPlayer(otherPlayer)
	}

	g.Chat.AddMessage(util.NewSystemMessage("GAME", "You've connected to GoldNet Official. Good luck!"))
}

func ClientPlayerJoinedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var pj packets.UpdatePlayer
	err := proto.Unmarshal(data, &pj)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal player joined packet")
		return
	}

	newPlayer := &game.Player{
		ID:       pj.Id,
		X:        pj.X,
		Username: pj.Username,
		HP:       pj.Hp,
		Hostile:  pj.Hostile,
	}
	gs.AddPlayer(newPlayer)

	g.Chat.AddMessage(util.NewSystemMessage("GAME", "A player has appeared!"))
}

func ClientPlayerDisconnectedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var pd packets.PlayerDisconnected
	err := proto.Unmarshal(data, &pd)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal player disconnected packet")
		return
	}

	g.Chat.AddMessage(util.NewSystemMessage("GAME", "A player has vanished!"))
	gs.RemovePlayer(pd.Id)
}
