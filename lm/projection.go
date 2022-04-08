package lm

import "math"

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

func TopDownProjection() Matrix3 {
	alpha := math.Pi / 6

	cosAlpha := float32(math.Cos(alpha))
	sinAlpha := float32(math.Sin(alpha))

	rotateX := NewMatrix3([9]float32{
		1, 0, 0,
		0, cosAlpha, sinAlpha,
		0, -sinAlpha, cosAlpha,
	})

	return rotateX
}
