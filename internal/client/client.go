package client

import (
	"bufio"
	"github.com/jessehorne/goldnet/internal/client/handlers"
	"github.com/jessehorne/goldnet/internal/game"
	"log"
	"net"
	"os"
)

type Client struct {
	Logger    *log.Logger
	Conn      net.Conn
	Done      chan struct{}
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
		Logger: log.New(os.Stdout, "[GoldNet] (Client) ", log.Ldate|log.Ltime),
		Conn:   conn,
		Reader: reader,
	}, nil
}

func (c *Client) Listen() {
	handler := handlers.NewPacketHandler(c.GameState)

	for {
		select {
		case <-c.Done:
			c.Close()
			return
		default:
			res, err := c.Reader.ReadBytes('\n')
			if err != nil {
				c.Logger.Println("[READ ERROR]", err.Error())
				continue
			}
			handler.Handle(c.Conn, res)
		}
	}
}

func (c *Client) Close() {
	c.Logger.Println("shutting down")
	c.Conn.Close()
}
