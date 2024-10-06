package client

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/client/handlers"
	"github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/util"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Client struct {
	Logger          *log.Logger
	Conn            net.Conn
	Reader          *bufio.Reader
	GameState       *game.GameState
	GUI             *gui.GUI
	HaveDrawnChunks bool
}

func NewClient() (*Client, error) {
	conn, err := net.Dial("tcp", ":5555")
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)

	c := &Client{
		Logger:    log.New(os.Stdout, "[GoldNet] (Client) ", log.Ldate|log.Ltime),
		Conn:      conn,
		Reader:    reader,
		GameState: game.NewGameState(),
	}

	g := gui.NewGUI(c.HandleInput)
	c.GUI = g

	go func(c *Client) {
		c.Loop()
	}(c)

	return c, nil
}

func (c *Client) Loop() {
	for {
		time.After(100 * time.Millisecond)
		// draw chunks for the first time if applicable
		if len(c.GameState.Chunks) > 0 && !c.HaveDrawnChunks {
			c.HaveDrawnChunks = true
			c.RedrawChunks()
			c.GUI.Chat.AddMessage(fmt.Sprintf("drawing %d chunks for the first time", len(c.GUI.World.Chunks)))
			continue
		}

		// redraw chunks if the player is in a different chunk
		playerID, ok := c.GameState.IntStore["playerID"]
		if !ok {
			continue
		}
		p := c.GameState.Players[playerID]
		chunkX := p.X / 8
		chunkY := p.Y / 8
		if chunkX != p.OldChunkX || chunkY != p.OldChunkY {
			p.OldChunkX = chunkX
			p.OldChunkY = chunkY
			c.RedrawChunks()
		}
	}
}

func (c *Client) RedrawChunks() {
	playerID, ok := c.GameState.IntStore["playerID"]
	if !ok {
		return
	}
	p := c.GameState.Players[playerID]

	chunks := c.GameState.GetChunksAroundPlayer(p)
	c.GUI.Chat.AddMessage(fmt.Sprintf("test: %d | p.ID = %d", len(chunks), p.ID))
	c.GUI.World.UpdateChunks(chunks)
	//gs.Logger.Println("received chunk", chunk.X, chunk.Y)
	c.GUI.Chat.AddMessage(util.NewSystemMessage("GAME", fmt.Sprintf("redrawing %d chunks", len(c.GUI.World.Chunks))))
}

func (c *Client) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	if c.GUI.World.Focused {
		playerID, ok := c.GameState.IntStore["playerID"]
		if !ok {
			return event
		}
		p := c.GameState.Players[playerID]

		switch event.Rune() {
		case 'a':
			p.X--
			c.GUI.World.OffsetX++
		case 'd':
			p.X++
			c.GUI.World.OffsetX--
		case 'w':
			p.Y--
			c.GUI.World.OffsetY++
		case 's':
			p.Y++
			c.GUI.World.OffsetY--
		}
	}
	return event
}

func (c *Client) Listen() {
	handler := handlers.NewPacketHandler(c.GameState)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	c.Conn.Write(packets.BuildUserJoinPacket())

	go func() {
		for {
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
			handler.Handle(c.GUI, c.Conn, data)
		}
	}()

	select {
	case <-done:
		c.Conn.Write(packets.BuildUserLeavePacket())
		c.Close()
		return
	}
}

func (c *Client) Close() {
	c.Logger.Println("shutting down")
	c.Conn.Close()
}
