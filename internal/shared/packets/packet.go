package packets

const (
	PacketUserJoin byte = iota
	PacketUserLeave
	PacketAction
	PacketSendMessage
	PacketChunks
	PacketPlayerJoined
	PacketPlayerSelfJoined
	PacketPlayerDisconnected
	PacketPlayerMoved
	PacketUpdateSelfPlayer
	PacketUpdateZombie
)

const (
	// Actions
	ActionMoveUp byte = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
)

var (
	MovementActions = []byte{ActionMoveUp, ActionMoveDown, ActionMoveLeft, ActionMoveRight}
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

func IsMovementAction(a byte) bool {
	for _, k := range MovementActions {
		if a == k {
			return true
		}
	}
	return false
}
