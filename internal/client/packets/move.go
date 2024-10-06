package packets

import (
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildUserMovePacket(action byte) []byte {
	p := util.Int64ToBytes(2)
	p = append(p, packets.PacketAction)
	p = append(p, action)
	return p
}
