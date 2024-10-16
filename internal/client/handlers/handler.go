package handlers

import (
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared/packets"
)

type Handler func(g *gui.GUI, gs *game.GameState, conn net.Conn, p []byte)

type PacketHandler struct {
	GameState *game.GameState
	Handlers  map[byte]Handler
}

func NewPacketHandler(gs *game.GameState) *PacketHandler {
	return &PacketHandler{
		GameState: gs,
		Handlers: map[byte]Handler{
			packets.PacketPlayerJoined:        ClientPlayerJoinedHandler,
			packets.PacketPlayerSelfJoined:    ClientPlayerSelfJoinedHandler,
			packets.PacketPlayerDisconnected:  ClientPlayerDisconnectedHandler,
			packets.PacketPlayerMoved:         ClientPlayerMovedHandler,
			packets.PacketPlayerToggleHostile: ClientPlayerToggleHostileHandler,
			packets.PacketChunks:              ClientChunksHandler,
			packets.PacketSendMessage:         ClientMessageHandler,
			packets.PacketUpdateSelfPlayer:    ClientUpdateSelfPlayerHandler,
			packets.PacketUpdateZombie:        ClientUpdateZombieHandler,
		},
	}
}

func (h *PacketHandler) Handle(g *gui.GUI, conn net.Conn, data []byte) {
	raw := packets.NewRawPacket(data)
	handler, ok := h.Handlers[raw.Type]
	if ok {
		handler(g, h.GameState, conn, raw.Data)
	}
}
