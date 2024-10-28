package handlers

import (
	"fmt"
	"net"

	"github.com/jessehorne/goldnet/internal/shared"
	packets "github.com/jessehorne/goldnet/packets/dist"

	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/client/handlers/components"
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
			shared.PacketUpdatePosition:     components.ClientUpdatePositionHandler,
			shared.PacketUpdateSprite:       components.ClientUpdateSpriteHandler,
			shared.PacketUpdatePlayer:       components.ClientUpdatePlayerHandler,
			shared.PacketUpdateZombie:       components.ClientUpdateZombieHandler,
			shared.PacketRemoveZombie:       components.ClientRemoveZombieHandler,
		},
	}
}

func (h *PacketHandler) Handle(g *gui.GUI, conn net.Conn, msg *packets.Raw, data []byte) {
	handler, ok := h.Handlers[msg.Type]
	if ok {
		handler(g, h.GameState, conn, data)
	} else {
		msg := fmt.Sprintf("[DEBUG ERROR] Received a message I don't know how to handle %d", msg.Type)
		g.Chat.AddMessage(msg)
	}
}
