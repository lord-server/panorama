package main

type Vector2 struct {
	X, Y, Z float32
}

func NewVector2(x, y float32) Vector2 {
	return Vector2{
		X: x,
		Y: y,
	}
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
