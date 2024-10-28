package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/game/components"
	"github.com/jessehorne/goldnet/internal/util"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
)

func ClientPlayerSelfJoinedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var sj packets.SelfJoin
	err := proto.Unmarshal(data, &sj)
	if err != nil {
		gs.Logger.Println(err)
		return
	}

	p := components.NewPlayer(components.EntityId(sj.Self.Id), sj.Self.Inventory, nil)
	gs.AddPlayer(p)
	gs.SetIntStore("playerID", sj.Self.Id)
	g.Sidebar.UpdatePlayerStats(p)

	g.Chat.AddMessage(util.NewSystemMessage("GAME", "You've connected to GoldNet Official. Good luck!"))
}

func ClientPlayerJoinedHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var pj packets.UpdatePlayer
	err := proto.Unmarshal(data, &pj)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal player joined packet")
		return
	}

	newPlayer := &components.Player{
		ID:       components.EntityId(pj.Id),
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
	gs.RemovePlayer(components.EntityId(pd.Id))
}
