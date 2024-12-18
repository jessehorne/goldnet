package shared

const (
	PacketUserJoin int32 = iota
	PacketUserLeave
	PacketAction
	PacketSendMessage
	PacketChunks
	PacketPlayerJoined
	PacketPlayerSelfJoined
	PacketPlayerDisconnected
	PacketUpdatePosition
	PacketUpdateSprite
	PacketUpdatePlayer
	PacketUpdateZombie
	PacketRemoveZombie
	PacketUseItem
	PacketSetHostile
)

const (
	// Actions
	ActionMoveUp int32 = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionToggleHostile
)

var (
	MovementActions = []int32{ActionMoveUp, ActionMoveDown, ActionMoveLeft, ActionMoveRight}
)

type RawPacket struct {
	Type byte
	Data []byte
}

func NewRawPacket(data []byte) *RawPacket {
	if len(data) == 0 {
		return nil
	}

	return &RawPacket{
		Type: data[0],
		Data: data[1:],
	}
}

func (p *RawPacket) ToBytes() []byte {
	return append([]byte{p.Type}, p.Data...)
}

func IsMovementAction(a int32) bool {
	for _, k := range MovementActions {
		if a == k {
			return true
		}
	}
	return false
}
