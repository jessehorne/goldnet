package inventory

import "github.com/jessehorne/goldnet/internal/util"

type Note struct {
	ID         int64
	Name       string
	Text       string
	ObjectType byte
	Callback   func()
	Quantity   int64
}

func NewNoteFromBytes(data []byte) *Note {
	counter := 0

	id := util.BytesToInt64(data[counter : counter+8])
	counter += 8

	nameLen := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	name := string(data[counter : counter+int(nameLen)])
	counter += int(nameLen)

	textLen := util.BytesToInt64(data[counter : counter+8])
	counter += 8
	text := string(data[counter : counter+int(textLen)])
	counter += int(textLen)

	quantity := util.BytesToInt64(data[counter : counter+8])

	return &Note{
		ID:         id,
		Name:       name,
		Text:       text,
		ObjectType: ObjectNote,
		Callback:   nil,
		Quantity:   quantity,
	}
}

func NewNote(qty int64, name, text string) *Note {
	return &Note{
		ID:         NextItemCounter(),
		Name:       name,
		Text:       text,
		ObjectType: ObjectNote,
		Quantity:   qty,
	}
}

// ToBytes returns the serialized version of the note.
/*
	===
	Packet
	===
	ObjectType (1 byte)
	ID (8 bytes)
	Length of Text (8 bytes)
	Text Data (X bytes)
	Quantity (8 bytes)
*/
func (n *Note) ToBytes() []byte {
	data := []byte{n.ObjectType}
	data = append(data, util.Int64ToBytes(n.ID)...)
	data = append(data, util.Int64ToBytes(int64(len(n.Name)))...)
	data = append(data, []byte(n.Name)...)
	data = append(data, util.Int64ToBytes(int64(len(n.Text)))...)
	data = append(data, []byte(n.Text)...)
	data = append(data, util.Int64ToBytes(n.Quantity)...)
	return append(util.Int64ToBytes(int64(len(data))), data...)
}

func (n *Note) GetName() string {
	return n.Name
}

func (n *Note) GetObjectType() byte {
	return n.ObjectType
}

func (n *Note) GetQuantity() int64 {
	return n.Quantity
}

func (n *Note) SetUseCallback(f func()) {
	n.Callback = f
}

func (n *Note) TriggerUse() {
	n.Callback()
}

func (n *Note) GetID() int64 {
	return n.ID
}
