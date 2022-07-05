package lm

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
	return Vec3(x, y, z)
}
