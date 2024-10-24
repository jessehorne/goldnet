package components

import (
	"time"

	"github.com/jessehorne/goldnet/internal/util"
)

const ZOMBIE_FOLLOW_RANGE int64 = 25

type ZombieComponent struct {
	ID                EntityId
	X                 int64
	Y                 int64
	HP                int64
	Damage            int64
	GoldDropAmt       int64
	FollowingPlayerId EntityId // -1 if not following anyone

	LastAttackTime time.Time
}

func NewZombieComponent(entityId EntityId, x, y int64) *ZombieComponent {
	mod := util.Distance(x, y, 0, 0) / 500
	return &ZombieComponent{
		ID:                entityId,
		X:                 x,
		Y:                 y,
		HP:                10 + mod,
		Damage:            5 + mod,
		GoldDropAmt:       5 + mod,
		FollowingPlayerId: -1,
		LastAttackTime:    time.Now(),
	}
}
