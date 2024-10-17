package packets

import "github.com/jessehorne/goldnet/internal/util"

func BuildUpdateZombiePacket(zombieData []byte) []byte {
	p := util.Int64ToBytes(int64(len(zombieData)) + 1)
	p = append(p, PacketUpdateZombie)
	p = append(p, zombieData...)
	return p
}

func BuildRemoveZombiePacket(zombieID int64) []byte {
	id_bytes := util.Int64ToBytes(zombieID)

	p := util.Int64ToBytes(int64(len(id_bytes)) + 1)
	p = append(p, PacketRemoveZombie)
	p = append(p, id_bytes...)
	return p
}
