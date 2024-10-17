package inventory

var (
	itemCounter int64 = 0
)

func NextItemCounter() int64 {
	itemCounter++
	return itemCounter
}

type Item interface {
	GetName() string
	GetID() int64
	GetObjectType() byte
	GetQuantity() int64
	ToBytes() []byte
	SetUseCallback(func())
	TriggerUse()
}
