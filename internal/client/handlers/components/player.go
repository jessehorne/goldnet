package components

import (
	"fmt"
	"net"

	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/game/components"
)

func ClientUpdatePlayerHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var up packets.UpdatePlayer
	err := proto.Unmarshal(data, &up)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal update player packet")
		return
	}

	msg := fmt.Sprintf("Handling update to player %d", up.Id)
	g.Chat.AddMessage(msg)

	p := components.NewPlayer(components.EntityId(up.Id), up.GetInventory(), nil)
	gs.PlayerComponents[components.EntityId(up.Id)] = p

	// gs.MovePlayer(p.ID, p.X, p.Y)

	currentPlayerID, exists := gs.GetIntStore("playerID")
	if exists {
		if currentPlayerID == up.Id {
			g.Sidebar.UpdatePlayerStats(p)
			// g.World.OffsetX = 50 + -int(p.X)
			// g.World.OffsetY = 13 + -int(p.Y)
		}
	}

}
