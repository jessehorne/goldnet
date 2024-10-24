package components

import (
	"fmt"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/game/components"
	packets "github.com/jessehorne/goldnet/packets/dist/components"
	"google.golang.org/protobuf/proto"
)

func ClientUpdateSpriteHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var up packets.UpdateSprite
	err := proto.Unmarshal(data, &up)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal update sprite packet")
		return
	}

	msg := fmt.Sprintf("Handling sprite update for player %d", up.EntityId)
	g.Chat.AddMessage(msg)

	gs.Mutex.Lock()
	gs.SpriteComponents[components.EntityId(up.EntityId)] = components.NewSpriteComponent(
		up.Character,
		tcell.Color(up.Foreground),
		tcell.Color(up.Background),
	)
	gs.Mutex.Unlock()
}
