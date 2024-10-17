package packets

import (
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildSendMessagePacket(msg string) []byte {
	msgBytes := []byte(msg)
	p := util.Int64ToBytes(1 + int64(len(msgBytes)))
	p = append(p, packets.PacketSendMessage)
	p = append(p, msgBytes...)
	return p
}

func ParseSendMessagePacket(data []byte) string {
	msg := string(data)
	return msg
}
