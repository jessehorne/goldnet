package game

import (
	"github.com/jessehorne/goldnet/internal/util"
)

const (
	CHUNK_W int64 = 8
	CHUNK_H int64 = 8
)

type Chunk struct {
	X     int64
	Y     int64
	Stack [CHUNK_H][CHUNK_W][]byte
}

func NewChunk(x, y int64) *Chunk {
	return &Chunk{
		X:     x,
		Y:     y,
		Stack: [CHUNK_H][CHUNK_W][]byte{},
	}
}

func (c *Chunk) FillPerlin() {
	for y := int64(0); y < CHUNK_H; y++ {
		for x := int64(0); x < CHUNK_W; x++ {
			coordsX := (c.X * CHUNK_W) + x
			coordsY := (c.Y * CHUNK_H) + y
			height := util.PerlinGetDataAtCoords(coordsX, coordsY)

			// the first byte of the stack is the height (0 - 255, 0 being lowest)
			c.Stack[y][x] = []byte{height}

			// handle water
			if height < 80 {
				c.Stack[y][x] = append(c.Stack[y][x], '~')
				continue
			}

			// handle sand
			if height >= 80 && height < 100 {
				c.Stack[y][x] = append(c.Stack[y][x], '.')
				continue
			}

			// handle lower and higher grass (100-150, 220-255)
			//isLowerGrass := height >= 100 && height < 150
			//isHigherGrass := height >= 220
			//if isLowerGrass || isHigherGrass {
			//	haveGrass := util.RandomIntBetween(0, 16) == 1
			//	if haveGrass {
			//		c.Stack[y][x] = append(c.Stack[y][x], ' ')
			//	}
			//	continue
			//}

			// handle middle-grass (where trees can grow)
			isMiddleGrass := height >= 150 && height < 220
			if isMiddleGrass {
				haveTree := util.RandomIntBetween(0, 8) == 1
				if haveTree {
					c.Stack[y][x] = append(c.Stack[y][x], 'T')
				}
				continue
			}
		}
	}
}

func (c *Chunk) ToBytes() []byte {
	var data []byte
	data = append(data, util.Int64ToBytes(c.X)...)
	data = append(data, util.Int64ToBytes(c.Y)...)
	for y := int64(0); y < CHUNK_H; y++ {
		for x := int64(0); x < CHUNK_W; x++ {
			stackSize := byte(len(c.Stack[y][x]))
			data = append(data, stackSize)
			for i := byte(0); i < stackSize; i++ {
				data = append(data, c.Stack[y][x][i])
			}
		}
	}
	return data
}

func ParseChunksFromBytes(data []byte) []*Chunk {
	chunkCount := util.BytesToInt64(data[0:8])
	counter := 8
	var chunks []*Chunk
	for cc := int64(0); cc < chunkCount; cc++ {
		chunkX := util.BytesToInt64(data[counter : counter+8])
		counter += 8
		chunkY := util.BytesToInt64(data[counter : counter+8])
		counter += 8
		chunkBytes := [CHUNK_H][CHUNK_W][]byte{}
		for y := int64(0); y < CHUNK_H; y++ {
			for x := int64(0); x < CHUNK_W; x++ {
				stackSize := data[counter]
				counter++

				for i := byte(0); i < stackSize; i++ {
					chunkBytes[y][x] = append(chunkBytes[y][x], data[counter])
					counter++
				}
			}
		}
		newChunk := &Chunk{
			X:     chunkX,
			Y:     chunkY,
			Stack: chunkBytes,
		}
		chunks = append(chunks, newChunk)
	}
	return chunks
}

func (c *Chunk) GetTopBlock(x, y int64) byte {
	stackSize := len(c.Stack[y][x])
	if stackSize == 1 {
		return ' '
	}
	return c.Stack[y][x][len(c.Stack[y][x])-1]
}
