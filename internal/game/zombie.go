package game

import (
	"github.com/jessehorne/goldnet/internal/util"
)

var (
	zombieCounter int64 = 0
)

type Zombie struct {
	ID          int64
	X           int64
	Y           int64
	HP          int64
	Damage      int64
	GoldDropAmt int64
}

func NewZombie(x, y int64) *Zombie {
	zombieCounter++
	mod := util.Distance(x, y, 0, 0) / 500
	return &Zombie{
		ID:          zombieCounter,
		X:           x,
		Y:           y,
		HP:          10 + mod,
		Damage:      5 + mod,
		GoldDropAmt: 5 + mod,
	}
}

func (z *Zombie) ToBytes() []byte {
	p := util.Int64ToBytes(z.ID)
	p = append(p, util.Int64ToBytes(z.X)...)
	p = append(p, util.Int64ToBytes(z.Y)...)
	p = append(p, util.Int64ToBytes(z.HP)...)
	p = append(p, util.Int64ToBytes(z.Damage)...)
	p = append(p, util.Int64ToBytes(z.GoldDropAmt)...)
	return p
}
