package game

import (
	"github.com/jessehorne/goldnet/internal/util"
)

const (
	CHUNK_W int64 = 8
	CHUNK_H int64 = 8
)

const (
	AboveNothing byte = iota
	AboveTree
)

type Chunk struct {
	X     int64
	Y     int64
	Below [CHUNK_H][CHUNK_W]byte
	Above [CHUNK_H][CHUNK_W]byte
}

func NewChunk(x, y int64) *Chunk {
	return &Chunk{
		X:     x,
		Y:     y,
		Below: [CHUNK_H][CHUNK_W]byte{},
		Above: [CHUNK_H][CHUNK_W]byte{},
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
			c.Below[y][x] = block
		}
	}
}

func (c *Chunk) FillPerlin() {
	for y := int64(0); y < CHUNK_H; y++ {
		for x := int64(0); x < CHUNK_W; x++ {
			coordsX := (c.X * CHUNK_W) + x
			coordsY := (c.Y * CHUNK_H) + y
			below := util.PerlinGetDataAtCoords(coordsX, coordsY)
			if below > 150 && below < 220 {
				haveTree := util.RandomIntBetween(0, 8)
				if haveTree == 1 {
					c.Above[y][x] = AboveTree
				} else {
					c.Above[y][x] = AboveNothing
				}
			}
			c.Below[y][x] = below
		}
	}
}

func (c *Chunk) ToBytes() []byte {
	var data []byte
	data = append(data, util.Int64ToBytes(c.X)...)
	data = append(data, util.Int64ToBytes(c.Y)...)
	for y := int64(0); y < CHUNK_H; y++ {
		for x := int64(0); x < CHUNK_W; x++ {
			data = append(data, c.Below[y][x])
		}
	}
	for y := int64(0); y < CHUNK_H; y++ {
		for x := int64(0); x < CHUNK_W; x++ {
			data = append(data, c.Above[y][x])
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
		chunkAboveBytes := [CHUNK_H][CHUNK_W]byte{}
		for y := int64(0); y < CHUNK_H; y++ {
			for x := int64(0); x < CHUNK_W; x++ {
				chunkBytes[y][x] = data[counter]
				counter++
			}
		}
		for y := int64(0); y < CHUNK_H; y++ {
			for x := int64(0); x < CHUNK_W; x++ {
				chunkAboveBytes[y][x] = data[counter]
				counter++
			}
		}
		newChunk := &Chunk{
			X:     chunkX,
			Y:     chunkY,
			Below: chunkBytes,
			Above: chunkAboveBytes,
		}
		chunks = append(chunks, newChunk)
	}
	return chunks
}

func (c *Chunk) GetBelowBlock(x, y int64) byte {
	if y > int64(len(c.Below)-1) {
		return 0
	}
	if x > int64(len(c.Below[y])-1) {
		return 0
	}
	return c.Below[y][x]
}

func (c *Chunk) GetAboveBlock(x, y int64) byte {
	if y > int64(len(c.Below)-1) {
		return 0
	}
	if x > int64(len(c.Below[y])-1) {
		return 0
	}
	return c.Above[y][x]
}
