package handlers

import (
	"github.com/jessehorne/goldnet/internal/shared"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"net"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/game"
)

type Handler func(g *gui.GUI, gs *game.GameState, conn net.Conn, p []byte)

type PacketHandler struct {
	GameState *game.GameState
	Handlers  map[int32]Handler
}

func NewPacketHandler(gs *game.GameState) *PacketHandler {
	return &PacketHandler{
		GameState: gs,
		Handlers: map[int32]Handler{
			shared.PacketPlayerJoined:       ClientPlayerJoinedHandler,
			shared.PacketPlayerSelfJoined:   ClientPlayerSelfJoinedHandler,
			shared.PacketPlayerDisconnected: ClientPlayerDisconnectedHandler,
			shared.PacketSetHostile:         ClientPlayerToggleHostileHandler,
			shared.PacketChunks:             ClientChunksHandler,
			shared.PacketSendMessage:        ClientMessageHandler,
			shared.PacketUpdatePlayer:       ClientUpdatePlayerHandler,
			shared.PacketUpdateZombie:       ClientUpdateZombieHandler,
			shared.PacketRemoveZombie:       ClientRemoveZombieHandler,
		},
	}
}

func (h *PacketHandler) Handle(g *gui.GUI, conn net.Conn, msg *packets.Raw, data []byte) {
	handler, ok := h.Handlers[msg.Type]
	if ok {
		handler(g, h.GameState, conn, data)
	}
}
