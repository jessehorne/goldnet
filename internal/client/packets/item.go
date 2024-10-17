package packets

import (
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildUseItemPacket(itemID int64) []byte {
	p := util.Int64ToBytes(1 + 8)
	p = append(p, packets.PacketUseItem)
	p = append(p, util.Int64ToBytes(itemID)...)
	return p
}
