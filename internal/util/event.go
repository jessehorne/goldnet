package util

import "slices"

func IsRuneMovementKey(r rune) bool {
	return slices.Contains([]rune{'a', 's', 'w', 'd'}, r)
}
