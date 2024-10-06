package packets

import "github.com/jessehorne/goldnet/internal/util"

func BuildPlayerJoinedPacket(playerID int64) []byte {
	p := util.Int64ToBytes(9) // total size should be 8 bytes (not counting this) because the playerID is 8 bytes
	p = append(p, PacketPlayerJoined)
	p = append(p, util.Int64ToBytes(playerID)...)
	return p
}

func BuildPlayerSelfJoinedPacket(playerID, x, y int64) []byte {
	p := util.Int64ToBytes(8*3 + 1) // total size should be 8 bytes (not counting this) because the playerID is 8 bytes
	p = append(p, PacketPlayerSelfJoined)
	p = append(p, util.Int64ToBytes(playerID)...)
	p = append(p, util.Int64ToBytes(x)...)
	p = append(p, util.Int64ToBytes(y)...)
	return p
}

func ParsePlayerSelfJoinedPacket(data []byte) (int64, int64, int64) {
	playerID := util.BytesToInt64(data[0:8])
	x := util.BytesToInt64(data[8:16])
	y := util.BytesToInt64(data[16:24])
	return playerID, x, y
}

func BuildPlayerDisconnectedPacket(playerID int64) []byte {
	p := util.Int64ToBytes(9) // total size should be 8 bytes (not counting this) because the playerID is 8 bytes
	p = append(p, PacketPlayerDisconnected)
	p = append(p, util.Int64ToBytes(playerID)...)
	return p
}
