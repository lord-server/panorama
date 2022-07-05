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

func (lhs Vector3) PowScalar(power float32) Vector3 {
	x := float32(math.Pow(float64(lhs.X), float64(power)))
	y := float32(math.Pow(float64(lhs.Y), float64(power)))
	z := float32(math.Pow(float64(lhs.Z), float64(power)))
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

func (lhs Vector3) RotateXY(angle float32) Vector3 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	return Vec3(lhs.X*cos-lhs.Y*sin, lhs.X*sin+lhs.Y*cos, lhs.Z)
}

func (lhs Vector3) RotateXZ(angle float32) Vector3 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	return Vec3(lhs.X*cos-lhs.Z*sin, lhs.Y, lhs.X*sin+lhs.Z*cos)
}

func (lhs Vector3) RotateYZ(angle float32) Vector3 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	return Vec3(lhs.X, lhs.Y*cos-lhs.Z*sin, lhs.Y*sin+lhs.Z*cos)
}
