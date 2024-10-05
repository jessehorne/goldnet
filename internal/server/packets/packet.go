package packets

const (
	// From Client

	PacketUserJoin byte = iota
	PacketUserLeave
	PacketInputKey
	PacketSendMessage
)

const (
	// To Client

	PacketChunk byte = iota
	PacketPlayerJoined
	PacketPlayerDisconnected
	PacketPlayerMoved
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
