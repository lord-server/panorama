package lm

import (
	"math"
)

func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

type Vector2 struct {
	X, Y float32
}

func NewVector2(x, y float32) Vector2 {
	return Vector2{
		X: x,
		Y: y,
	}
}

func (lhs Vector2) Add(rhs Vector2) Vector2 {
	x := lhs.X + rhs.X
	y := lhs.Y + rhs.Y
	return NewVector2(x, y)
}

func (lhs Vector2) Sub(rhs Vector2) Vector2 {
	x := lhs.X - rhs.X
	y := lhs.Y - rhs.Y
	return NewVector2(x, y)
}

func (lhs Vector2) Mul(rhs Vector2) Vector2 {
	x := lhs.X * rhs.X
	y := lhs.Y * rhs.Y
	return NewVector2(x, y)
}

func (lhs Vector2) MulScalar(rhs float32) Vector2 {
	x := lhs.X * rhs
	y := lhs.Y * rhs
	return NewVector2(x, y)
}

func (lhs Vector2) Min(rhs Vector2) Vector2 {
	x := minFloat32(lhs.X, rhs.X)
	y := minFloat32(lhs.Y, rhs.Y)
	return NewVector2(x, y)
}

func (lhs Vector2) Max(rhs Vector2) Vector2 {
	x := maxFloat32(lhs.X, rhs.X)
	y := maxFloat32(lhs.Y, rhs.Y)
	return NewVector2(x, y)
}

type Vector3 struct {
	X, Y, Z float32
}

func NewVector3(x, y, z float32) Vector3 {
	return Vector3{
		X: x,
		Y: y,
		Z: z,
	}
}

func (lhs Vector3) Add(rhs Vector3) Vector3 {
	x := lhs.X + rhs.X
	y := lhs.Y + rhs.Y
	z := lhs.Z + rhs.Z
	return NewVector3(x, y, z)
}

func (lhs Vector3) MulScalar(rhs float32) Vector3 {
	x := lhs.X * rhs
	y := lhs.Y * rhs
	z := lhs.Z * rhs
	return NewVector3(x, y, z)
}

func (lhs Vector3) DivScalar(rhs float32) Vector3 {
	reciprocal := 1 / rhs
	x := lhs.X * reciprocal
	y := lhs.Y * reciprocal
	z := lhs.Z * reciprocal
	return NewVector3(x, y, z)
}

func (lhs Vector3) Cross(rhs Vector3) Vector3 {
	x := lhs.Y*rhs.Z - lhs.Z*rhs.Y
	y := lhs.Z*rhs.X - lhs.X*rhs.Z
	z := lhs.X*rhs.Y - lhs.Y*rhs.X
	return NewVector3(x, y, z)
}

func (lhs Vector3) Dot(rhs Vector3) float32 {
	return lhs.X*rhs.X + lhs.Y*rhs.Y + lhs.Z*rhs.Z
}

func (lhs Vector3) Length() float32 {
	return float32(math.Sqrt(float64(lhs.Dot(lhs))))
}

func (lhs Vector3) Normalize() Vector3 {
	return lhs.DivScalar(lhs.Length())
}

func (lhs Vector3) ClampScalar(min, max float32) Vector3 {
	x := Clamp(lhs.X, min, max)
	y := Clamp(lhs.Y, min, max)
	z := Clamp(lhs.Z, min, max)
	return NewVector3(x, y, z)
}

func (v Vector3) XY() Vector2 {
	return NewVector2(v.X, v.Y)
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

type Matrix3 struct {
	m [9]float32
}

func NewMatrix3(m [9]float32) Matrix3 {
	return Matrix3{
		m: m,
	}
}

func (lhs *Matrix3) Mul(rhs *Matrix3) Matrix3 {
	return NewMatrix3([9]float32{
		lhs.m[0]*rhs.m[0] + lhs.m[1]*rhs.m[3] + lhs.m[2]*rhs.m[6],
		lhs.m[0]*rhs.m[1] + lhs.m[1]*rhs.m[4] + lhs.m[2]*rhs.m[7],
		lhs.m[0]*rhs.m[2] + lhs.m[1]*rhs.m[5] + lhs.m[2]*rhs.m[8],
		lhs.m[3]*rhs.m[0] + lhs.m[4]*rhs.m[3] + lhs.m[5]*rhs.m[6],
		lhs.m[3]*rhs.m[1] + lhs.m[4]*rhs.m[4] + lhs.m[5]*rhs.m[7],
		lhs.m[3]*rhs.m[2] + lhs.m[4]*rhs.m[5] + lhs.m[5]*rhs.m[8],
		lhs.m[6]*rhs.m[0] + lhs.m[7]*rhs.m[3] + lhs.m[8]*rhs.m[6],
		lhs.m[6]*rhs.m[1] + lhs.m[7]*rhs.m[4] + lhs.m[8]*rhs.m[7],
		lhs.m[6]*rhs.m[2] + lhs.m[7]*rhs.m[5] + lhs.m[8]*rhs.m[8],
	})
}

func (lhs *Matrix3) MulVec(rhs Vector3) Vector3 {
	x := lhs.m[0]*rhs.X + lhs.m[1]*rhs.Y + lhs.m[2]*rhs.Z
	y := lhs.m[3]*rhs.X + lhs.m[4]*rhs.Y + lhs.m[5]*rhs.Z
	z := lhs.m[6]*rhs.X + lhs.m[7]*rhs.Y + lhs.m[8]*rhs.Z
	return NewVector3(x, y, z)
}

func DimetricProjection() Matrix3 {
	alpha := math.Pi / 6
	beta := math.Pi / 4

	cosAlpha := float32(math.Cos(alpha))
	sinAlpha := float32(math.Sin(alpha))

	cosBeta := float32(math.Cos(beta))
	sinBeta := float32(math.Sin(beta))

	rotateX := NewMatrix3([9]float32{
		1, 0, 0,
		0, cosAlpha, sinAlpha,
		0, -sinAlpha, cosAlpha,
	})

	rotateY := NewMatrix3([9]float32{
		cosBeta, 0, -sinBeta,
		0, 1, 0,
		sinBeta, 0, cosBeta,
	})

	return rotateX.Mul(&rotateY)
}
