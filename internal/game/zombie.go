package game

import (
	"time"

	"github.com/jessehorne/goldnet/internal/game/components"
	"github.com/jessehorne/goldnet/internal/util"
)

var (
	zombieCounter components.EntityId = 0
)

const ZOMBIE_FOLLOW_RANGE int64 = 25

type Zombie struct {
	ID                components.EntityId
	X                 int64
	Y                 int64
	HP                int64
	Damage            int64
	GoldDropAmt       int64
	FollowingPlayerId components.EntityId // -1 if not following anyone

	LastAttackTime time.Time
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
		LastAttackTime:    time.Now(),
	}
}

func (z *Zombie) ToBytes() []byte {
	p := util.Int64ToBytes(int64(z.ID))
	p = append(p, util.Int64ToBytes(z.HP)...)
	p = append(p, util.Int64ToBytes(z.Damage)...)
	p = append(p, util.Int64ToBytes(z.GoldDropAmt)...)
	p = append(p, util.Int64ToBytes(int64(z.FollowingPlayerId))...)
	return p
}
