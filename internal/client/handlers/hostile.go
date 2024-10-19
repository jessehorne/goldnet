package handlers

import (
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
)

func ClientPlayerToggleHostileHandler(g *gui.GUI, gs *game.GameState, conn net.Conn, data []byte) {
	var shp packets.SetHostile
	err := proto.Unmarshal(data, &shp)
	if err != nil {
		gs.Logger.Println("couldn't unmarshal set hostile packet")
		return
	}

	gs.Mutex.Lock()
	gs.Players[shp.PlayerID].Hostile = shp.Hostile
	gs.Mutex.Unlock()
}
