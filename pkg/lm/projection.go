package lm

import "math"

func DimetricProjection() Matrix3 {
	alpha := math.Pi / 6
	beta := math.Pi / 4

	cosAlpha := math.Cos(alpha)
	sinAlpha := math.Sin(alpha)

	cosBeta := math.Cos(beta)
	sinBeta := math.Sin(beta)

	rotateX := NewMatrix3([9]float64{
		1, 0, 0,
		0, cosAlpha, sinAlpha,
		0, -sinAlpha, cosAlpha,
	})

	rotateY := NewMatrix3([9]float64{
		cosBeta, 0, -sinBeta,
		0, 1, 0,
		sinBeta, 0, cosBeta,
	})

	return rotateX.Mul(&rotateY)
}

func TopDownProjection() Matrix3 {
	alpha := math.Pi / 2

	cosAlpha := math.Cos(alpha)
	sinAlpha := math.Sin(alpha)
	rotateY := NewMatrix3([9]float64{
		1, 0, 0,
		0, cosAlpha, sinAlpha,
		0, -sinAlpha, cosAlpha,
	})

	scale := NewMatrix3([9]float64{
		math.Sqrt2, 0, 0,
		0, math.Sqrt2, 0,
		0, 0, math.Sqrt2,
	})

	return scale.Mul(&rotateY)
}
