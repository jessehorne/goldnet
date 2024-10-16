package game

import (
	"github.com/jessehorne/goldnet/internal/util"
)

var (
	zombieCounter int64 = 0
)

const ZOMBIE_FOLLOW_RANGE int64 = 25

type Zombie struct {
	ID                int64
	X                 int64
	Y                 int64
	HP                int64
	Damage            int64
	GoldDropAmt       int64
	FollowingPlayerId int64 // -1 if not following anyone
}

func NewZombie(x, y int64) *Zombie {
	zombieCounter++
	mod := util.Distance(x, y, 0, 0) / 500
	return &Zombie{
		ID:                zombieCounter,
		X:                 x,
		Y:                 y,
		HP:                10 + mod,
		Damage:            5 + mod,
		GoldDropAmt:       5 + mod,
		FollowingPlayerId: -1,
	}
}

func (z *Zombie) ToBytes() []byte {
	p := util.Int64ToBytes(z.ID)
	p = append(p, util.Int64ToBytes(z.X)...)
	p = append(p, util.Int64ToBytes(z.Y)...)
	p = append(p, util.Int64ToBytes(z.HP)...)
	p = append(p, util.Int64ToBytes(z.Damage)...)
	p = append(p, util.Int64ToBytes(z.GoldDropAmt)...)
	p = append(p, util.Int64ToBytes(z.FollowingPlayerId)...)
	return p
}
