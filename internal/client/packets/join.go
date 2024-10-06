package packets

import (
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildUserJoinPacket() []byte {
	p := util.Int64ToBytes(1) // only 1 byte because we're just sending the type (not counting the size of the data)
	p = append(p, packets.PacketUserJoin)
	return p
}

func BuildUserLeavePacket() []byte {
	p := util.Int64ToBytes(1)
	p = append(p, packets.PacketUserLeave)
	return p
}
