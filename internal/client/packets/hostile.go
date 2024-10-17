package packets

import (
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildUserToggleHostilePacket(hostile bool) []byte {
	p := util.Int64ToBytes(2)
	p = append(p, packets.PacketAction)
	p = append(p, packets.ActionToggleHostile)
	return p
}
