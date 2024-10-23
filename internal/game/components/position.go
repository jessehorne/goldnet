package components

type PositionComponent struct {
	X int64
	Y int64
}

func NewPositionComponent(x, y int64) *PositionComponent {
	return &PositionComponent{
		X: x,
		Y: y,
	}
}
