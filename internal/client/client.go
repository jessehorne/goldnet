package client

import (
	"bufio"
	"github.com/jessehorne/goldnet/internal/client/handlers"
	"github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/util"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Client struct {
	Logger    *log.Logger
	Conn      net.Conn
	Reader    *bufio.Reader
	GameState *game.GameState
}

func NewClient() (*Client, error) {
	conn, err := net.Dial("tcp", ":5555")
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)

	return &Client{
		Logger:    log.New(os.Stdout, "[GoldNet] (Client) ", log.Ldate|log.Ltime),
		Conn:      conn,
		Reader:    reader,
		GameState: game.NewGameState(),
	}, nil
}

func (c *Client) Listen() {
	handler := handlers.NewPacketHandler(c.GameState)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-done:
			c.Conn.Write(packets.BuildUserLeavePacket())
			c.Close()
			return
		default:
			// first 8 bytes (int64) is how large this packet is in bytes
			var sizeBytes []byte
			for i := 0; i < 8; i++ {
				b, err := c.Reader.ReadByte()
				if err != nil {
					continue
				}
				sizeBytes = append(sizeBytes, b)
			}

			if len(sizeBytes) != 8 {
				continue
			}

			// convert size to int64
			size := util.BytesToInt64(sizeBytes)

			// read that many bytes which is the packet
			var data []byte
			for i := int64(0); i < size; i++ {
				b, err := c.Reader.ReadByte()
				if err != nil {
					continue
				}
				data = append(data, b)
			}
			handler.Handle(c.Conn, data)
		}
	}
}

func (c *Client) Close() {
	c.Logger.Println("shutting down")
	c.Conn.Close()
}
