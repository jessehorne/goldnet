package inventory

import (
	"github.com/jessehorne/goldnet/internal/util"
)

type Inventory struct {
	Items []Item
}

func NewInventory(data []byte) *Inventory {
	newInv := &Inventory{
		Items: []Item{},
	}

	// parse items if there are any
	if len(data) > 0 {
		numberOfItems := util.BytesToInt64(data[0:8])

		counter := 8
		for i := int64(0); i < numberOfItems; i++ {
			sizeOfData := util.BytesToInt64(data[counter : counter+8])
			counter += 8
			objType := data[counter]
			counter++

			if objType == ObjectNote {
				newNote := NewNoteFromBytes(data[counter : counter+int(sizeOfData)-1]) // subtracting one because the size includes the obj type which we don't need
				counter += int(sizeOfData) - 1
				newInv.Items = append(newInv.Items, newNote)
			}
		}
	}

	return newInv
}

func (i *Inventory) AddItem(item Item) {
	i.Items = append(i.Items, item)
}

func (i *Inventory) ToBytes() []byte {
	var data []byte
	for _, item := range i.Items {
		data = append(data, item.ToBytes()...)
	}
	return append(util.Int64ToBytes(int64(len(i.Items))), data...)
}
