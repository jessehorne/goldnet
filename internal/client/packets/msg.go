package packets

import (
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
)

func BuildMessagePacket(playerID int64, data string) []byte {
	dataBytes := []byte(data)
	p := util.Int64ToBytes(1 + 8 + int64(len(dataBytes)))
	p = append(p, packets.PacketSendMessage)
	p = append(p, util.Int64ToBytes(playerID)...)
	p = append(p, dataBytes...)
	return p
}

func ParseMessagePacket(data []byte) (int64, string) {
	playerID := util.BytesToInt64(data[0:8])
	msg := string(data[8:])
	return playerID, msg
}

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
