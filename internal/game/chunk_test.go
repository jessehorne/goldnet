package game

import "testing"

func Test_NewChunk(t *testing.T) {
	nc := NewChunk(0, 0)
	if nc.Data[5][5] != 0 {
		t.Error("error creating chunk")
	}
}
