package game

const (
	CHUNK_W int = 8
	CHUNK_H int = 8
)

type Chunk struct {
	X    int64
	Y    int64
	Data [CHUNK_H][CHUNK_W]byte
}

func NewChunk(x, y int64) *Chunk {
	return &Chunk{
		X:    x,
		Y:    y,
		Data: [CHUNK_H][CHUNK_W]byte{},
	}
}
