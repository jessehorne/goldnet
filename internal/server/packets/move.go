package packets

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildMovePacket(playerID, x, y int64) []byte {
	p := []byte{PacketPlayerMoved}
	p = append(p, util.Int64ToBytes(playerID)...)
	p = append(p, util.Int64ToBytes(x)...)
	p = append(p, util.Int64ToBytes(y)...)
	p = append(p, '\n')
	return p
}

func ParseMovePacket(data []byte) (int64, int64, int64) {
	fmt.Println("SIZE:", len(data))
	playerID := util.BytesToInt64(data[0:8])
	x := util.BytesToInt64(data[8:16])
	y := util.BytesToInt64(data[16:24])
	return playerID, x, y
}
