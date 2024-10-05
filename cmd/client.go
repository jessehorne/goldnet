package main

import (
	"bufio"
	"fmt"
	"github.com/jessehorne/goldnet/internal/server/packets"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":5555")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		conn.Write([]byte{packets.PacketUserJoin, byte('\n')})
	}()

	reader := bufio.NewReader(conn)

	for {
		select {
		case <-done:
			time.Sleep(1 * time.Second)
			conn.Write([]byte{packets.PacketUserLeave, byte('\n')})
			fmt.Println("shutting down")
			return
		default:
			res, err := reader.ReadBytes(byte('\n'))
			if err != nil {
				panic(err)
			}
			if res[0] == packets.PacketPlayerJoined {
				fmt.Println("another player joined")
			} else if res[0] == packets.PacketPlayerDisconnected {
				fmt.Println("a player disconnected")
			}
		}
	}
}
