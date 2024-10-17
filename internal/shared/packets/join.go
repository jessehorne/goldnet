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

func ParsePlayerJoinedPacket(data []byte) (int64, int64, int64) {
	playerID := util.BytesToInt64(data[0:8])
	x := util.BytesToInt64(data[8:16])
	y := util.BytesToInt64(data[16:24])
	return playerID, x, y
}

func BuildPlayerSelfJoinedPacket(playerID, x, y int64, players []byte, inventory []byte) []byte {
	var p []byte
	p = append(p, PacketPlayerSelfJoined)
	p = append(p, util.Int64ToBytes(playerID)...)
	p = append(p, util.Int64ToBytes(x)...)
	p = append(p, util.Int64ToBytes(y)...)

	// add all data on other players
	p = append(p, util.Int64ToBytes(int64(len(players)))...)
	p = append(p, players...)

	// add all data on players own inventory
	p = append(p, util.Int64ToBytes(int64(len(inventory)))...)
	p = append(p, inventory...)

	// the first 8 bytes should be the int64 size of the packet
	p = append(util.Int64ToBytes(int64(len(p))), p...)

	return p
}

func ParsePlayerSelfJoinedPacket(data []byte) (int64, int64, int64, []byte, []byte) {
	counter := 0

	playerID := util.BytesToInt64(data[counter : counter+8])
	counter += 8

	x := util.BytesToInt64(data[counter : counter+8])
	counter += 8

	y := util.BytesToInt64(data[counter : counter+8])
	counter += 8

	playerDataLen := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	otherPlayers := data[counter : counter+int(playerDataLen)]
	counter += int(playerDataLen)

	inventoryDataLen := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	inventoryData := data[counter : counter+int(inventoryDataLen)]
	counter += int(inventoryDataLen)

	return playerID, x, y, otherPlayers, inventoryData
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
