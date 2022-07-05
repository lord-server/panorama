package lm

import "math"

func Abs(value float32) float32 {
	return float32(math.Abs(float64(value)))
}

func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func Radians(degrees float32) float32 {
	return degrees * math.Pi / 180.0
}
