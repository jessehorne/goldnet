package packets

import "github.com/jessehorne/goldnet/internal/util"

type ClientJoinPacket struct{}

func BuildPlayerJoinedPacket(playerID int64) []byte {
	p := []byte{PacketPlayerJoined}
	p = append(p, util.Int64ToBytes(playerID)...)
	p = append(p, '\n')
	return p
}

func ParsePlayerJoinedPacket(data []byte) int64 {
	return util.BytesToInt64(data)
}
