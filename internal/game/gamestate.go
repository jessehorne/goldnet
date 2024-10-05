package game

import (
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

func (gs *GameState) AddPlayer(playerID int64, c net.Conn) {
	gs.Players[playerID] = NewPlayer(playerID, 0, 0, c)
}

func (gs *GameState) RemovePlayer(playerID int64) {
	gs.Players[playerID].Conn.Close()
	gs.Players[playerID] = nil
}
