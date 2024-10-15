package packets

import "github.com/jessehorne/goldnet/internal/util"

func BuildNewZombiePacket(zombieData []byte) []byte {
	p := util.Int64ToBytes(int64(len(zombieData)) + 1)
	p = append(p, PacketNewZombie)
	p = append(p, zombieData...)
	return p
}

func BuildUpdateZombiePacket(zombieData []byte) []byte {
	p := util.Int64ToBytes(int64(len(zombieData)) + 1)
	p = append(p, PacketUpdateZombie)
	p = append(p, zombieData...)
	return p
}
