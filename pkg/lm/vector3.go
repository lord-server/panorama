package lm

import "math"

type Vector3 struct {
	X, Y, Z float64
}

func Vec3(x, y, z float64) Vector3 {
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

func (lhs Vector3) MulScalar(rhs float64) Vector3 {
	x := lhs.X * rhs
	y := lhs.Y * rhs
	z := lhs.Z * rhs

	return Vec3(x, y, z)
}

func (lhs Vector3) DivScalar(rhs float64) Vector3 {
	reciprocal := 1 / rhs
	x := lhs.X * reciprocal
	y := lhs.Y * reciprocal
	z := lhs.Z * reciprocal

	return Vec3(x, y, z)
}

func (lhs Vector3) PowScalar(power float64) Vector3 {
	x := math.Pow(lhs.X, power)
	y := math.Pow(lhs.Y, power)
	z := math.Pow(lhs.Z, power)

	return Vec3(x, y, z)
}

func (lhs Vector3) Cross(rhs Vector3) Vector3 {
	x := lhs.Y*rhs.Z - lhs.Z*rhs.Y
	y := lhs.Z*rhs.X - lhs.X*rhs.Z
	z := lhs.X*rhs.Y - lhs.Y*rhs.X

	return Vec3(x, y, z)
}

func (lhs Vector3) Dot(rhs Vector3) float64 {
	return lhs.X*rhs.X + lhs.Y*rhs.Y + lhs.Z*rhs.Z
}

func (lhs Vector3) Length() float64 {
	return math.Sqrt(lhs.Dot(lhs))
}

func (lhs Vector3) Normalize() Vector3 {
	return lhs.DivScalar(lhs.Length())
}

func (lhs Vector3) ClampScalar(min, max float64) Vector3 {
	x := Clamp(lhs.X, min, max)
	y := Clamp(lhs.Y, min, max)
	z := Clamp(lhs.Z, min, max)

	return Vec3(x, y, z)
}

func (lhs Vector3) XY() Vector2 {
	return Vec2(lhs.X, lhs.Y)
}

func (lhs Vector3) MaxComponent() float64 {
	return math.Max(lhs.X, math.Max(lhs.Y, lhs.Z))
}

func (lhs Vector3) RotateXY(angle float64) Vector3 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	return Vec3(lhs.X*cos-lhs.Y*sin, lhs.X*sin+lhs.Y*cos, lhs.Z)
}

func (lhs Vector3) RotateXZ(angle float64) Vector3 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	return Vec3(lhs.X*cos-lhs.Z*sin, lhs.Y, lhs.X*sin+lhs.Z*cos)
}

func (lhs Vector3) RotateYZ(angle float64) Vector3 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	return Vec3(lhs.X, lhs.Y*cos-lhs.Z*sin, lhs.Y*sin+lhs.Z*cos)
}
