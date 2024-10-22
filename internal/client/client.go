package client

import (
	"bufio"
	packets "github.com/jessehorne/goldnet/packets/dist"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/goldnet/internal/client/gui"
	"github.com/jessehorne/goldnet/internal/client/handlers"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/jessehorne/goldnet/internal/shared"
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
					useItem := &packets.UseItem{
						Type: shared.PacketUseItem,
						Id:   int64(c.GUI.Sidebar.InventoryCursor),
					}
					useItemData, useItemErr := proto.Marshal(useItem)
					if useItemErr != nil {
						c.GameState.Logger.Println("couldn't marshal use item packet")
						return event
					}
					util.Send(c.Conn, useItemData)
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
				p.Hostile = !p.Hostile
				hp := &packets.SetHostile{
					Type:     shared.PacketSetHostile,
					PlayerID: p.ID,
					Hostile:  p.Hostile,
				}
				hpData, hpErr := proto.Marshal(hp)
				if hpErr != nil {
					c.GameState.Logger.Println("couldn't marshal move set hostile packet")
					return event
				}
				util.Send(c.Conn, hpData)

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

				um := &packets.Move{
					Type:   shared.PacketAction,
					Action: shared.ActionMoveLeft,
				}
				umData, umErr := proto.Marshal(um)
				if umErr != nil {
					c.GameState.Logger.Println("couldn't marshal move left")
					return event
				}
				util.Send(c.Conn, umData)

				p.LastMovementTime = time.Now()
			}
		case 'd':
			if canMove {
				p.X++
				c.GUI.World.OffsetX--
				c.GUI.World.OldOffsetX = c.GUI.World.OffsetX
				p.OldChunkX = p.X / 8

				um := &packets.Move{
					Type:   shared.PacketAction,
					Action: shared.ActionMoveRight,
				}
				umData, umErr := proto.Marshal(um)
				if umErr != nil {
					c.GameState.Logger.Println("couldn't marshal move right")
					return event
				}
				util.Send(c.Conn, umData)

				p.LastMovementTime = time.Now()
			}
		case 'w':
			if canMove {
				p.Y--
				c.GUI.World.OffsetY++
				c.GUI.World.OldOffsetY = c.GUI.World.OffsetY
				p.OldChunkY = p.Y / 8

				um := &packets.Move{
					Type:   shared.PacketAction,
					Action: shared.ActionMoveUp,
				}
				umData, umErr := proto.Marshal(um)
				if umErr != nil {
					c.GameState.Logger.Println("couldn't marshal move up")
					return event
				}
				util.Send(c.Conn, umData)

				p.LastMovementTime = time.Now()
			}
		case 's':
			if canMove {
				p.Y++
				c.GUI.World.OffsetY--
				c.GUI.World.OldOffsetY = c.GUI.World.OffsetY
				p.OldChunkY = p.Y / 8

				um := &packets.Move{
					Type:   shared.PacketAction,
					Action: shared.ActionMoveDown,
				}
				umData, umErr := proto.Marshal(um)
				if umErr != nil {
					c.GameState.Logger.Println("couldn't marshal move down")
					return event
				}
				util.Send(c.Conn, umData)
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
				sm := &packets.Message{
					Type: shared.PacketSendMessage,
					Data: msg,
				}
				smData, smErr := proto.Marshal(sm)
				if smErr != nil {
					c.GameState.Logger.Println(smErr)
					return event
				}
				util.Send(c.Conn, smData)
			}
			return event
		}
	}
	return event
}

func (c *Client) Listen() {
	handler := handlers.NewPacketHandler(c.GameState)

	uj := &packets.Join{
		Ptype: shared.PacketUserJoin,
	}
	ujData, ujErr := proto.Marshal(uj)
	if ujErr != nil {
		c.GameState.Logger.Println("couldn't marshal initial user join packet")
		return
	}
	util.Send(c.Conn, ujData)

	go func() {
		for {
			lenBytes := make([]byte, 8)
			_, err := io.ReadFull(c.Reader, lenBytes)
			if err != nil {
				continue
			}

			msgLen := util.BytesToInt64(lenBytes)
			msgBytes := make([]byte, msgLen)
			io.ReadFull(c.Reader, msgBytes)

			msg := packets.Raw{}
			err = proto.Unmarshal(msgBytes, &msg)
			if err != nil {
				continue
			}

			handler.Handle(c.GUI, c.Conn, &msg, msgBytes)
		}
	}()
}

func (c *Client) Close() {
	ul := &packets.Leave{
		Type: shared.PacketUserLeave,
	}
	ulData, ulErr := proto.Marshal(ul)
	if ulErr != nil {
		c.GameState.Logger.Println("couldn't marshal user leave packet")
		return
	}
	util.Send(c.Conn, ulData)
	c.Conn.Close()
	c.App.Stop()
}
