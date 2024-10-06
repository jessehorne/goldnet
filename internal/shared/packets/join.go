package packets

import "github.com/jessehorne/goldnet/internal/util"

func BuildPlayerJoinedPacket(playerID int64) []byte {
	p := util.Int64ToBytes(9) // total size should be 8 bytes (not counting this) because the playerID is 8 bytes
	p = append(p, PacketPlayerJoined)
	p = append(p, util.Int64ToBytes(playerID)...)
	return p
}

func BuildPlayerDisconnectedPacket(playerID int64) []byte {
	p := util.Int64ToBytes(9) // total size should be 8 bytes (not counting this) because the playerID is 8 bytes
	p = append(p, PacketPlayerDisconnected)
	p = append(p, util.Int64ToBytes(playerID)...)
	return p
}
