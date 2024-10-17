package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/util"
	"net"
)

func ServerUseItemHandler(gs *game.GameState, playerID int64, conn net.Conn, data []byte) {
	p := gs.GetPlayer(playerID)
	if p == nil {
		return
	}
	itemID := util.BytesToInt64(data)
	if int64(len(p.Inventory.Items)) > itemID {
		p.Inventory.Items[itemID].TriggerUse()
	}
}
