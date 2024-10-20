package client

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/client/handlers"
	"github.com/jessehorne/goldnet/internal/client/packets"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared"
	sharedPackets "github.com/jessehorne/goldnet/internal/shared/packets"
	"github.com/jessehorne/goldnet/internal/util"
	"github.com/rivo/tview"
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

		// handle inventory
		if event.Key() == tcell.KeyPgUp {
			if c.GUI.Sidebar.InventoryCursor > 0 {
				c.GUI.Sidebar.InventoryCursor--
			} else {
				c.GUI.Sidebar.InventoryCursor = len(p.Inventory.Items) - 1
			}
			c.GUI.Sidebar.UpdatePlayerInventory(p)
			return nil
		} else if event.Key() == tcell.KeyPgDn {
			if c.GUI.Sidebar.InventoryCursor < len(p.Inventory.Items)-1 {
				c.GUI.Sidebar.InventoryCursor++
			} else {
				c.GUI.Sidebar.InventoryCursor = 0
			}
			c.GUI.Sidebar.UpdatePlayerInventory(p)
			return nil
		} else if event.Key() == tcell.KeyEnter {
			n, _ := c.GUI.Sidebar.Pages.GetFrontPage()
			if n == "inventory" {
				if len(p.Inventory.Items) > 0 {
					c.Conn.Write(packets.BuildUseItemPacket(int64(c.GUI.Sidebar.InventoryCursor)))
				}
			}
		}

		// handle chat
		if event.Rune() == 't' {
			c.GUI.World.Focused = false
			c.GUI.Input.Focused = true
			c.App.SetFocus(c.GUI.Input.Root)
			return nil
		}

		// determine if movement or something else
		isMovement := util.IsRuneMovementKey(event.Rune())
		if !isMovement {
			switch event.Rune() {
			// Toggle hostile mode
			case 'e':
				c.Conn.Write(packets.BuildUserToggleHostilePacket(p.Hostile))

			// Handle sidebar input
			case 'S':
				c.GUI.Sidebar.UpdatePlayerStats(p)
				c.GUI.Sidebar.SetActiveTab("stats")
			case 'I':
				c.GUI.Sidebar.UpdatePlayerInventory(p)
				c.GUI.Sidebar.SetActiveTab("inventory")
			}
			return nil
		}

		// it is movement, so handle movement
		mod := (1 / float64(p.Speed)) * 1000

		b := c.GameState.GetTerrainAtCoords(p.X, p.Y)
		if b == shared.TerrainWater {
			mod = mod * 4
		}

		canMoveAt := p.LastMovementTime.Add(time.Duration(mod) * time.Millisecond)
		canMove := true
		if time.Now().Before(canMoveAt) {
			canMove = false
		}

		switch event.Rune() {
		case 'a':
			if canMove {
				p.X--
				c.GUI.World.OldOffsetX = c.GUI.World.OffsetX
				c.GUI.World.OffsetX++
				p.OldChunkX = p.X / 8
				c.Conn.Write(packets.BuildUserMovePacket(sharedPackets.ActionMoveLeft))
				p.LastMovementTime = time.Now()
			}
		case 'd':
			if canMove {
				p.X++
				c.GUI.World.OffsetX--
				c.GUI.World.OldOffsetX = c.GUI.World.OffsetX
				p.OldChunkX = p.X / 8
				c.Conn.Write(packets.BuildUserMovePacket(sharedPackets.ActionMoveRight))
				p.LastMovementTime = time.Now()
			}
		case 'w':
			if canMove {
				p.Y--
				c.GUI.World.OffsetY++
				c.GUI.World.OldOffsetY = c.GUI.World.OffsetY
				p.OldChunkY = p.Y / 8
				c.Conn.Write(packets.BuildUserMovePacket(sharedPackets.ActionMoveUp))
				p.LastMovementTime = time.Now()
			}
		case 's':
			if canMove {
				p.Y++
				c.GUI.World.OffsetY--
				c.GUI.World.OldOffsetY = c.GUI.World.OffsetY
				p.OldChunkY = p.Y / 8
				c.Conn.Write(packets.BuildUserMovePacket(sharedPackets.ActionMoveDown))
				p.LastMovementTime = time.Now()
			}
		}
	} else if c.GUI.Input.Focused {
		switch event.Key() {
		case tcell.KeyESC:
			c.GUI.World.Focused = true
			c.GUI.Input.Focused = false
			c.App.SetFocus(c.GUI.World.Root)
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
