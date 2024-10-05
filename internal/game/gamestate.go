package game

import (
	"github.com/jessehorne/goldnet/internal/server/packets"
	"github.com/jessehorne/goldnet/internal/util"
	"log"
	"net"
	"os"
)

type GameState struct {
	Players map[int64]*Player
	Logger  *log.Logger
}

func NewGameState() *GameState {
	return &GameState{
		Logger:  log.New(os.Stdout, "[GoldNet] (GameState) ", log.Ldate|log.Ltime),
		Players: map[int64]*Player{},
	}
}

func (gs *GameState) GetPlayer(playerID int64) *Player {
	p, ok := gs.Players[playerID]
	if !ok {
		return nil
	}
	return p
}

func (gs *GameState) AddPlayer(playerID int64, c net.Conn) {
	gs.Players[playerID] = NewPlayer(playerID, 0, 0, c)
}

func (gs *GameState) RemovePlayer(playerID int64) {
	gs.Players[playerID].Conn.Close()
	gs.Players[playerID] = nil
}

func (gs *GameState) HandlePlayerAction(playerID int64, action byte) {
	p := gs.GetPlayer(playerID)
	if p != nil {
		p.Action(action)

		if packets.IsMovementAction(action) {
			// send movement to any players nearby
			movePacket := packets.BuildMovePacket(playerID, p.X, p.Y)
			for _, o := range gs.Players {
				if o != nil {
					if o.ID != playerID {
						if util.Distance(o.X, o.Y, p.X, p.Y) < 100 {
							p.Conn.Write(movePacket)
						}
					}
				}
			}
		}
	}
}
