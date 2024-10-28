package gui

import (
	"fmt"

	"github.com/jessehorne/goldnet/internal/game/components"
	"github.com/rivo/tview"
)

type Sidebar struct {
	Root            *tview.Grid
	Pages           *tview.Pages
	Nav             *tview.TextView
	Stats           *tview.TextView
	Inventory       *tview.TextView
	InventoryCursor int
}

func NewSidebar() *Sidebar {
	grid := tview.NewGrid()
	grid.SetBorder(false)

	nav := tview.NewTextView()
	nav.SetBorder(true)
	nav.SetDynamicColors(true)
	nav.SetText("[white](S)tats    [grey](I)nventory")

	pages := tview.NewPages()
	pages.SetBorder(true)

	stats := tview.NewTextView()
	stats.SetBorder(false)
	stats.SetDynamicColors(true)
	pages.AddPage("stats", stats, true, true)

	inv := tview.NewTextView()
	inv.SetBorder(false)
	inv.SetDynamicColors(true)
	pages.AddPage("inventory", inv, true, false)

	pages.SwitchToPage("stats")
	pages.SetTitle("Player Stats")

	grid.AddItem(nav, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(pages, 1, 0, 6, 1, 0, 0, false)

	return &Sidebar{
		Root:      grid,
		Nav:       nav,
		Pages:     pages,
		Stats:     stats,
		Inventory: inv,
	}
}

func (s *Sidebar) SetActiveTab(name string) {
	s.Pages.SwitchToPage(name)
	var tmpl string
	if name == "stats" {
		s.Pages.SetTitle("Player Stats")
		tmpl = "[white](S)tats    [grey](I)nventory"
	} else if name == "inventory" {
		s.Pages.SetTitle("Player Inventory (PgUp / PgDown / Enter)")
		tmpl = "[grey](S)tats    [white](I)nventory"
	}
	s.Nav.SetText(tmpl)
}

func (s *Sidebar) UpdatePlayerStats(p *components.Player) {
	tmpl := `
Name: %s

[yellow]Gold: %d
[green]HP: %s%d
[blue]ST: %d
[white]
`
	var lowHealthText string
	if p.HP < 30 {
		lowHealthText = "[red]"
	}
	s.Stats.SetText(fmt.Sprintf(tmpl,
		p.Username,
		p.Gold,
		lowHealthText,
		p.HP,
		p.ST,
	))
}

func (s *Sidebar) UpdatePlayerInventory(p *components.Player) {
	tmpl := ``

	for index, item := range p.Inventory.Items {
		highlighted := index == s.InventoryCursor
		var color string
		if highlighted {
			color = "[white]"
		} else {
			color = "[grey]"
		}
		tmpl += fmt.Sprintf("%s%dx %s\n", color, item.GetQuantity(), item.GetName())
	}

	s.Inventory.SetText(tmpl)
}
