package packets

import (
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildSetHostilePacket(playerID int64, hostile bool) []byte {
	p := util.Int64ToBytes(8*2 + 1)
	p = append(p, PacketPlayerToggleHostile)
	p = append(p, util.Int64ToBytes(playerID)...)

	if hostile {
		p = append(p, util.Int64ToBytes(1)...)
	} else {
		p = append(p, util.Int64ToBytes(0)...)
	}
	return p
}

func ParseSetHostilePacket(data []byte) (int64, bool) {
	playerID := util.BytesToInt64(data[0:8])
	hostile := util.BytesToInt64(data[8:16]) == 1
	return playerID, hostile
}
