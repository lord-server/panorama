package lm

import "math"

type Vector2 struct {
	X, Y float64
}

func Vec2(x, y float64) Vector2 {
	return Vector2{
		X: x,
		Y: y,
	}
}

func (lhs Vector2) Add(rhs Vector2) Vector2 {
	x := lhs.X + rhs.X
	y := lhs.Y + rhs.Y
	return Vec2(x, y)
}

func (lhs Vector2) Sub(rhs Vector2) Vector2 {
	x := lhs.X - rhs.X
	y := lhs.Y - rhs.Y
	return Vec2(x, y)
}

func (lhs Vector2) Mul(rhs Vector2) Vector2 {
	x := lhs.X * rhs.X
	y := lhs.Y * rhs.Y
	return Vec2(x, y)
}

func (lhs Vector2) MulScalar(rhs float64) Vector2 {
	x := lhs.X * rhs
	y := lhs.Y * rhs
	return Vec2(x, y)
}

func (lhs Vector2) Min(rhs Vector2) Vector2 {
	x := math.Min(lhs.X, rhs.X)
	y := math.Min(lhs.Y, rhs.Y)
	return Vec2(x, y)
}

func (lhs Vector2) Max(rhs Vector2) Vector2 {
	x := math.Max(lhs.X, rhs.X)
	y := math.Max(lhs.Y, rhs.Y)
	return Vec2(x, y)
}
