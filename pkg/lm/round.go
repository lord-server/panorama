package lm

import "math"

// floorDiv returns the result of floor division. The difference compared to
// usual division is that floor division always rounds down instead of rounding
// towards zero.
func FloorDiv(a, b int) int {
	return int(math.Floor(float64(a) / float64(b)))
}
