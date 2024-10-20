package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"net"
)

type Handler func(gs *game.GameState, playerID int64, conn net.Conn, p []byte)

type PacketHandler struct {
	GameState *game.GameState
	Handlers  map[int32]Handler
}

func NewPacketHandler(gs *game.GameState) *PacketHandler {
	return &PacketHandler{
		GameState: gs,
		Handlers: map[int32]Handler{
			shared.PacketUserJoin:    ServerUserJoinHandler,
			shared.PacketUserLeave:   ServerUserDisconnectedHandler,
			shared.PacketAction:      ServerActionHandler,
			shared.PacketSendMessage: ServerMessageHandler,
			shared.PacketUseItem:     ServerUseItemHandler,
			shared.PacketSetHostile:  ServerSetHostileHandler,
		},
	}
}

func (h *PacketHandler) Handle(playerID int64, conn net.Conn, msg *packets.Raw, data []byte) {
	handler, ok := h.Handlers[msg.Type]
	if ok {
		handler(h.GameState, playerID, conn, data)
	}
}
