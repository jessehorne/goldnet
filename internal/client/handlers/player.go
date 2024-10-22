package handlers

import (
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
)

func ClientUpdatePlayerHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var up packets.UpdatePlayer
	err := proto.Unmarshal(data, &up)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal update player packet")
		return
	}
	p := gs.GetPlayer(up.Id)
	p.X = up.X
	p.Y = up.Y
	p.Gold = up.Gold
	p.HP = up.Hp
	p.ST = up.St
	p.Hostile = up.Hostile

	gs.MovePlayer(p.ID, p.X, p.Y)

	currentPlayerID, exists := gs.GetIntStore("playerID")
	if exists {
		if currentPlayerID == up.Id {
			g.Sidebar.UpdatePlayerStats(p)
			g.World.OffsetX = 50 + -int(p.X)
			g.World.OffsetY = 13 + -int(p.Y)
		}
	}

}
