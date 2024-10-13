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
	return &Sidebar{
		Root: tv,
	}
}

func (s *Sidebar) AttachPlayer(p *game.Player) {
	s.Player = p
}

func (s *Sidebar) UpdateText() {
	tmpl := `
%s

HP: %d
ST: %d

Inventory
---------

coming soon...
`
	s.Root.SetText(fmt.Sprintf(tmpl,
		s.Player.Username,
		s.Player.HP,
		s.Player.ST,
	))
}
