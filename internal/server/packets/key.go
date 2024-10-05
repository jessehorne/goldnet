package packets

import (
	"unicode/utf8"
)

type ClientKeyPacket struct {
	Key rune
}

func NewClientKeyPacket(raw *RawPacket) *ClientKeyPacket {
	r, _ := utf8.DecodeRune(raw.Data)
	return &ClientKeyPacket{
		Key: r,
	}
}
