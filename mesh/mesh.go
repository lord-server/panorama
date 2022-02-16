package mesh

import (
	"github.com/weqqr/panorama/lm"
)

type Vertex struct {
	Position lm.Vector3
	Texcoord lm.Vector2
	Normal   lm.Vector3
}

type Mesh struct {
	Vertices []Vertex
}

func NewMesh() Mesh {
	return Mesh{
		Vertices: []Vertex{},
	}
}

type Model struct {
	Meshes []Mesh
}

func NewModel() Model {
	return Model{
		Meshes: []Mesh{},
	}
}

func Cube() *Model {
	model := NewModel()
	yp := NewMesh()
	yp.Vertices = []Vertex{
		{Position: lm.Vec3(-0.5, +0.5, -0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(-0.5, +0.5, +0.5), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(+0.5, +0.5, +0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(-0.5, +0.5, -0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(+0.5, +0.5, -0.5), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(+0.5, +0.5, +0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
	}

	ym := NewMesh()
	ym.Vertices = []Vertex{
		{Position: lm.Vec3(-0.5, -0.5, -0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(-0.5, -0.5, +0.5), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(+0.5, -0.5, +0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(-0.5, -0.5, -0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(+0.5, -0.5, -0.5), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(+0.5, -0.5, +0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
	}

	xp := NewMesh()
	xp.Vertices = []Vertex{
		{Position: lm.Vec3(+0.5, -0.5, -0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(+0.5, -0.5, +0.5), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(+0.5, +0.5, +0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(+0.5, -0.5, -0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(+0.5, +0.5, -0.5), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(+0.5, +0.5, +0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
	}

	xm := NewMesh()
	xm.Vertices = []Vertex{
		{Position: lm.Vec3(-0.5, -0.5, -0.5), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(-0.5, -0.5, +0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(-0.5, +0.5, +0.5), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(-0.5, -0.5, -0.5), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(-0.5, +0.5, -0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(-0.5, +0.5, +0.5), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
	}

	zp := NewMesh()
	zp.Vertices = []Vertex{
		{Position: lm.Vec3(-0.5, -0.5, +0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(-0.5, +0.5, +0.5), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(+0.5, +0.5, +0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(-0.5, -0.5, +0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(+0.5, -0.5, +0.5), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(+0.5, +0.5, +0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
	}

	zm := NewMesh()
	zm.Vertices = []Vertex{
		{Position: lm.Vec3(-0.5, -0.5, -0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(-0.5, +0.5, -0.5), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(+0.5, +0.5, -0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(-0.5, -0.5, -0.5), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(+0.5, -0.5, -0.5), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(+0.5, +0.5, -0.5), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
	}

	model.Meshes = append(model.Meshes, yp)
	model.Meshes = append(model.Meshes, ym)
	model.Meshes = append(model.Meshes, xp)
	model.Meshes = append(model.Meshes, xm)
	model.Meshes = append(model.Meshes, zp)
	model.Meshes = append(model.Meshes, zm)

	return &model
}
