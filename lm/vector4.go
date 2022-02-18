package lm

type Vector4 struct {
	X, Y, Z, W float32
}

func Vec4(x, y, z, w float32) Vector4 {
	return Vector4{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}
}

func (lhs Vector4) MulScalar(rhs float32) Vector4 {
	x := lhs.X * rhs
	y := lhs.Y * rhs
	z := lhs.Z * rhs
	w := lhs.W * rhs
	return Vec4(x, y, z, w)
}

func (lhs Vector4) ClampScalar(min, max float32) Vector4 {
	x := Clamp(lhs.X, min, max)
	y := Clamp(lhs.Y, min, max)
	z := Clamp(lhs.Z, min, max)
	w := Clamp(lhs.W, min, max)
	return Vec4(x, y, z, w)
}

func (lhs Vector4) XYZ() Vector3 {
	return Vec3(lhs.X, lhs.Y, lhs.Z)
}
