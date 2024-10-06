package packets

type ClientMessagePacket struct{}

func NewClientMessagePacket(p *RawPacket) ClientMessagePacket {
	return ClientMessagePacket{}
}
