package main

import (
	"github.com/jessehorne/goldnet/internal/client"
	clientPackets "github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/shared/packets"
	"time"
)

func main() {
	c, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	go func() {
		c.Conn.Write(clientPackets.BuildUserJoinPacket())

		for {
			c.Conn.Write(clientPackets.BuildUserMovePacket(packets.ActionMoveLeft))
			time.Sleep(1 * time.Second)
		}
	}()

	c.Listen()
}
