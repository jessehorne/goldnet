package packets

import (
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildPlayerJoinedPacket(playerID, x, y int64) []byte {
	p := util.Int64ToBytes(8*3 + 1)
	p = append(p, PacketPlayerJoined)
	p = append(p, util.Int64ToBytes(playerID)...)
	p = append(p, util.Int64ToBytes(x)...)
	p = append(p, util.Int64ToBytes(y)...)
	return p
}

func ParsePlayerSelfJoinedPacket(data []byte) (int64, int64, int64, []byte) {
	playerID := util.BytesToInt64(data[0:8])
	x := util.BytesToInt64(data[8:16])
	y := util.BytesToInt64(data[16:24])
	otherPlayersData := data[24:]
	return playerID, x, y, otherPlayersData
}

func BuildPlayerSelfJoinedPacket(playerID, x, y int64, players []byte) []byte {
	var p = make([]byte, 8)
	p = append(p, PacketPlayerSelfJoined)
	p = append(p, util.Int64ToBytes(playerID)...)
	p = append(p, util.Int64ToBytes(x)...)
	p = append(p, util.Int64ToBytes(y)...)
	p = append(p, players...)
	l := util.Int64ToBytes(1 + 8*3 + int64(len(players)))
	for i := 0; i < 8; i++ {
		p[i] = l[i]
	}
	return p
}

func ParsePlayerJoinedPacket(data []byte) (int64, int64, int64) {
	playerID := util.BytesToInt64(data[0:8])
	x := util.BytesToInt64(data[8:16])
	y := util.BytesToInt64(data[16:24])
	return playerID, x, y
}

func BuildPlayerDisconnectedPacket(playerID int64) []byte {
	p := util.Int64ToBytes(9)
	p = append(p, PacketPlayerDisconnected)
	p = append(p, util.Int64ToBytes(playerID)...)
	return p
}

func ParsePlayerDisconnectedPacket(data []byte) int64 {
	playerID := util.BytesToInt64(data[0:8])
	return playerID
}
