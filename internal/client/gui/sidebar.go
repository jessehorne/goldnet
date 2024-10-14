package gui

import (
	"fmt"
	"github.com/jessehorne/goldnet/internal/game"
	"github.com/rivo/tview"
)

type Sidebar struct {
	Root   *tview.TextView
	Player *game.Player
}

func NewSidebar() *Sidebar {
	tv := tview.NewTextView()
	tv.SetBorder(false)
	tv.SetDynamicColors(true)
	return &Sidebar{
		Root: tv,
	}
}

func (s *Sidebar) AttachPlayer(p *game.Player) {
	s.Player = p
}

func (s *Sidebar) UpdateText() {
	tmpl := `
Player Name: %s

[yellow]Gold: %d
[green]HP: %s%d
[blue]ST: %d
[white]
Inventory
---------

coming soon...
`
	var lowHealthText string
	if s.Player.HP < 30 {
		lowHealthText = "[red]"
	}
	s.Root.SetText(fmt.Sprintf(tmpl,
		s.Player.Username,
		s.Player.Gold,
		lowHealthText,
		s.Player.HP,
		s.Player.ST,
	))
}
