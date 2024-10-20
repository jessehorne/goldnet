package util

import (
	"encoding/binary"
	"net"
)

func Int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func BytesToInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}

func Send(conn net.Conn, data []byte) {
	dataLen := Int64ToBytes(int64(len(data)))
	conn.Write(dataLen)
	conn.Write(data)
}
