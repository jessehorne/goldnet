package packets

import "github.com/jessehorne/goldnet/internal/util"

func BuildUpdateSelfPlayerPacket(data []byte) []byte {
	p := util.Int64ToBytes(1 + int64(len(data)))
	p = append(p, PacketUpdateSelfPlayer)
	p = append(p, data...)
	return p
}
