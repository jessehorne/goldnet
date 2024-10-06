package gui

import (
	"fmt"
	"github.com/rivo/tview"
)

type Chat struct {
	Root       *tview.TextView
	Focused    bool
	Messages   []string
	MaxHistory int
}

func NewChat() *Chat {
	tv := tview.NewTextView()
	tv.SetScrollable(true)
	tv.SetDynamicColors(true)

	return &Chat{
		Root:       tv,
		Messages:   []string{},
		MaxHistory: 1000,
	}
}

func (c *Chat) AddMessage(s string) {
	c.Messages = append(c.Messages, s)
	if len(c.Messages) > c.MaxHistory {
		c.Messages = c.Messages[1:]
	}

	all := ""
	for _, m := range c.Messages {
		all = fmt.Sprintf("%s%s\n", all, m)
	}

	c.Root.SetText(all)
	c.Root.ScrollToEnd()
}
