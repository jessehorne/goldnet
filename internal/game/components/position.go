package components

type Position struct {
	X int64
	Y int64
}

func NewPositionComponent(x, y int64) *Position {
	return &Position{
		X: x,
		Y: y,
	}
}
