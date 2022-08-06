package lm

import "math"

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
