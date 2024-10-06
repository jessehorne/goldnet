package packets

import (
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildMovePacket(playerID, x, y int64) []byte {
	p := util.Int64ToBytes(8*3 + 1)
	p = append(p, PacketPlayerMoved)
	p = append(p, util.Int64ToBytes(playerID)...)
	p = append(p, util.Int64ToBytes(x)...)
	p = append(p, util.Int64ToBytes(y)...)
	return p
}

func ParseMovePacket(data []byte) (int64, int64, int64) {
	playerID := util.BytesToInt64(data[0:8])
	x := util.BytesToInt64(data[8:16])
	y := util.BytesToInt64(data[16:24])
	return playerID, x, y
}
