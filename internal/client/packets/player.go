package packets

import (
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/util"
)

func ParsePlayerBytes(data []byte) *game.Player {
	usernameLen := util.BytesToInt64(data[0:8])
	var usernameData []byte
	counter := int64(8)
	for i := int64(0); i < usernameLen; i++ {
		usernameData = append(usernameData, data[counter])
	}
	username := string(usernameData)

	counter += usernameLen
	id := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	x := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	y := util.BytesToInt64(data[counter : counter+8])
	counter += 8

	gold := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	hp := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	st := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	hostileInt := util.BytesToInt64(data[counter : counter+8])
	counter += 8

	return &game.Player{
		ID:       id,
		Username: username,
		X:        x,
		Y:        y,
		Gold:     gold,
		HP:       hp,
		ST:       st,
		Hostile:  hostileInt == 1,
	}
}
