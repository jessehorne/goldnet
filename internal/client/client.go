package client

import (
	"bufio"
	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/client/handlers"
	"github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
	sharedPackets "github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
	"github.com/rivo/tview"
	"log"
	"net"
	"os"
)

type Client struct {
	Logger          *log.Logger
	Conn            net.Conn
	Reader          *bufio.Reader
	GameState       *game.GameState
	GUI             *gui.GUI
	HaveDrawnChunks bool
	App             *tview.Application
}

func NewClient(tv *tview.Application) (*Client, error) {
	conn, err := net.Dial("tcp", ":5555")
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	gs := game.NewGameState()
	c := &Client{
		Logger:    log.New(os.Stdout, "[GoldNet] (Client) ", log.Ldate|log.Ltime),
		Conn:      conn,
		Reader:    reader,
		GameState: gs,
		App:       tv,
	}

	g := gui.NewGUI(gs, c.HandleInput)
	c.GUI = g

	return c, nil
}

func (c *Client) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	if c.GUI.World.Focused {
		playerID, ok := c.GameState.GetIntStore("playerID")
		if !ok {
			return event
		}
		p := c.GameState.GetPlayer(playerID)
		if p == nil {
			return event
		}

		switch event.Rune() {
		case 'a':
			p.X--
			c.GUI.World.OffsetX++
			if p.OldChunkX != p.X/8 {
				p.OldChunkX = p.X / 8
			}
			c.Conn.Write(packets.BuildUserMovePacket(sharedPackets.ActionMoveLeft))
		case 'd':
			p.X++
			c.GUI.World.OffsetX--
			if p.OldChunkX != p.X/8 {
				p.OldChunkX = p.X / 8
			}
			c.Conn.Write(packets.BuildUserMovePacket(sharedPackets.ActionMoveRight))
		case 'w':
			p.Y--
			c.GUI.World.OffsetY++
			if p.OldChunkY != p.Y/8 {
				p.OldChunkY = p.Y / 8
			}
			c.Conn.Write(packets.BuildUserMovePacket(sharedPackets.ActionMoveUp))
		case 's':
			p.Y++
			c.GUI.World.OffsetY--
			if p.OldChunkY != p.Y/8 {
				p.OldChunkY = p.Y / 8
			}
			c.Conn.Write(packets.BuildUserMovePacket(sharedPackets.ActionMoveDown))
		case 't':
			if c.GUI.World.Focused {
				c.GUI.World.Focused = false
				c.GUI.Input.Focused = true
				i := c.GUI.Input.Root.GetFormItemByLabel("> ").(*tview.InputField)
				i.SetText("")
				c.App.SetFocus(c.GUI.Input.Root)
				return nil
			}
		}
	} else if c.GUI.Input.Focused {
		switch event.Key() {
		case tcell.KeyESC:
			c.GUI.World.Focused = true
			c.GUI.Input.Focused = false
			c.App.SetFocus(c.GUI.World.Root)
			i := c.GUI.Input.Root.GetFormItemByLabel("> ").(*tview.InputField)
			i.SetText("")
			return event
		case tcell.KeyEnter:
			c.GUI.World.Focused = true
			c.GUI.Input.Focused = false
			c.App.SetFocus(c.GUI.World.Root)
			i := c.GUI.Input.Root.GetFormItemByLabel("> ").(*tview.InputField)
			msg := i.GetText()
			i.SetText("")

			if len(msg) > 0 {
				c.Conn.Write(packets.BuildSendMessagePacket(msg))
			}
			return event
		}
	}
	return event
}

func (c *Client) Listen() {
	handler := handlers.NewPacketHandler(c.GameState)
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

			// read that many bytes which is the packet
			size := util.BytesToInt64(sizeBytes)
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
}

func (c *Client) Close() {
	c.Conn.Write(packets.BuildUserLeavePacket())
	c.Conn.Close()
	c.App.Stop()
}
