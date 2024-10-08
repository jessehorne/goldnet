package game

import (
	"github.com/jessehorne/goldnet/internal/util"
)

const (
	CHUNK_W int64 = 8
	CHUNK_H int64 = 8
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

func (c *Chunk) FillV1() {
	for y := int64(0); y < CHUNK_H; y++ {
		for x := int64(0); x < CHUNK_W; x++ {
			var block byte
			r := util.RandomIntBetween(0, 10)
			if r < 1 {
				block = ','
			} else {
				block = ' '
			}
			c.Data[y][x] = block
		}
	}
}

func (c *Chunk) ToBytes() []byte {
	var data []byte
	data = append(data, util.Int64ToBytes(c.X)...)
	data = append(data, util.Int64ToBytes(c.Y)...)
	for y := int64(0); y < CHUNK_H; y++ {
		for x := int64(0); x < CHUNK_W; x++ {
			data = append(data, c.Data[y][x])
		}
	}
	return data
}

func ParseChunksFromBytes(data []byte) []*Chunk {
	chunkCount := util.BytesToInt64(data[0:8])
	counter := 8
	var chunks []*Chunk
	for i := int64(0); i < chunkCount; i++ {
		chunkX := util.BytesToInt64(data[counter : counter+8])
		counter += 8
		chunkY := util.BytesToInt64(data[counter : counter+8])
		counter += 8
		chunkBytes := [CHUNK_H][CHUNK_W]byte{}
		for y := int64(0); y < CHUNK_H; y++ {
			for x := int64(0); x < CHUNK_W; x++ {
				chunkBytes[y][x] = data[counter]
				counter++
			}
		}
		newChunk := &Chunk{
			X:    chunkX,
			Y:    chunkY,
			Data: chunkBytes,
		}
		chunks = append(chunks, newChunk)
	}
	return chunks
}

func GetChunkFromCoords(x, y int64) (int64, int64) {
	return x / CHUNK_W, y / CHUNK_H
}
