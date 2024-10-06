package handlers

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"net"
)

type Handler func(gs *game.GameState, conn net.Conn, p []byte)

type PacketHandler struct {
	GameState *game.GameState
	Handlers  map[byte]Handler
}

func NewPacketHandler(gs *game.GameState) *PacketHandler {
	return &PacketHandler{
		GameState: gs,
		Handlers: map[byte]Handler{
			packets.PacketPlayerJoined:       ClientPlayerJoinedHandler,
			packets.PacketPlayerDisconnected: ClientPlayerDisconnectedHandler,
			packets.PacketPlayerMoved:        ClientPlayerMovedHandler,
			packets.PacketChunk:              ClientChunkHandler,
		},
	}
}

func (h *PacketHandler) Handle(conn net.Conn, data []byte) {
	raw := packets.NewRawPacket(data)
	handler, ok := h.Handlers[raw.Type]
	if ok {
		handler(h.GameState, conn, raw.Data)
	}
}
