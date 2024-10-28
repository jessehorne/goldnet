package components

import (
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/game/components"
	packets "github.com/jessehorne/goldnet/packets/dist/components"
	"google.golang.org/protobuf/proto"
)

func ClientUpdatePositionHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var up packets.UpdatePosition
	err := proto.Unmarshal(data, &up)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal update position packet")
		return
	}

	// msg := fmt.Sprintf("Handling position update for entity %d", up.EntityId)
	// g.Chat.AddMessage(msg)

	gs.Mutex.Lock()
	gs.PositionComponents[components.EntityId(up.EntityId)] = components.NewPositionComponent(up.X, up.Y)
	gs.Mutex.Unlock()
}
