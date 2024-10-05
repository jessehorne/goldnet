package packets

type ClientJoinPacket struct{}

func NewClientJoinPacket(p *RawPacket) ClientJoinPacket {
	return ClientJoinPacket{}
}
