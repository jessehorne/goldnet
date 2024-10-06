package packets

type ClientLeavePacket struct{}

func NewClientLeavePacket(p *RawPacket) ClientLeavePacket {
	return ClientLeavePacket{}
}
