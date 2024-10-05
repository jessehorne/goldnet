package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/server/packets"
	"net"
)

type Handler func(gs *game.GameState, playerID int64, conn net.Conn, p []byte)

type PacketHandler struct {
	GameState *game.GameState
	Handlers  map[byte]Handler
}

func NewPacketHandler(gs *game.GameState) *PacketHandler {
	return &PacketHandler{
		GameState: gs,
		Handlers: map[byte]Handler{
			packets.PacketUserJoin:    ClientUserJoinHandler,
			packets.PacketUserLeave:   ClientUserLeaveHandler,
			packets.PacketAction:      ClientActionHandler,
			packets.PacketSendMessage: ClientMessageHandler,
		},
	}
}

func (h *PacketHandler) Handle(playerID int64, conn net.Conn, data []byte) {
	raw := packets.NewRawPacket(data)
	handler, ok := h.Handlers[raw.Type]
	if ok {
		handler(h.GameState, playerID, conn, raw.Data)
	}
}
