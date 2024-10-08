package packets

import (
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildChunksPacket(chunkCount int64, data []byte) []byte {
	p := util.Int64ToBytes(1 + 8 + int64(len(data)))
	p = append(p, PacketChunks)
	p = append(p, util.Int64ToBytes(chunkCount)...)
	p = append(p, data...)
	return p
}
