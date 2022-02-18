package lm

import "math"

type Vector3 struct {
	X, Y, Z float32
}

func Vec3(x, y, z float32) Vector3 {
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
	return Vec3(x, y, z)
}

func (lhs Vector3) MulScalar(rhs float32) Vector3 {
	x := lhs.X * rhs
	y := lhs.Y * rhs
	z := lhs.Z * rhs
	return Vec3(x, y, z)
}

func (lhs Vector3) DivScalar(rhs float32) Vector3 {
	reciprocal := 1 / rhs
	x := lhs.X * reciprocal
	y := lhs.Y * reciprocal
	z := lhs.Z * reciprocal
	return Vec3(x, y, z)
}

func (lhs Vector3) Cross(rhs Vector3) Vector3 {
	x := lhs.Y*rhs.Z - lhs.Z*rhs.Y
	y := lhs.Z*rhs.X - lhs.X*rhs.Z
	z := lhs.X*rhs.Y - lhs.Y*rhs.X
	return Vec3(x, y, z)
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
	return Vec3(x, y, z)
}

func (lhs Vector3) XY() Vector2 {
	return Vec2(lhs.X, lhs.Y)
}

func (lhs Vector3) MaxComponent() float32 {
	return maxFloat32(lhs.X, maxFloat32(lhs.Y, lhs.Z))
}
