package mesh

import (
	"github.com/weqqr/panorama/pkg/lm"
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

func Cuboid(x1, y1, z1, x2, y2, z2 float32) []Mesh {
	yp := NewMesh()
	yp.Vertices = []Vertex{
		{Position: lm.Vec3(x1, y2, z1), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(x1, y2, z2), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(x2, y2, z2), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(x1, y2, z1), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(x2, y2, z1), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
		{Position: lm.Vec3(x2, y2, z2), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 1.0, 0.0)},
	}

	ym := NewMesh()
	ym.Vertices = []Vertex{
		{Position: lm.Vec3(x1, y1, z1), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(x1, y1, z2), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(x2, y1, z2), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(x1, y1, z1), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(x2, y1, z1), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
		{Position: lm.Vec3(x2, y1, z2), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, -1.0, 0.0)},
	}

	xp := NewMesh()
	xp.Vertices = []Vertex{
		{Position: lm.Vec3(x2, y1, z1), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x2, y1, z2), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x2, y2, z2), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x2, y1, z1), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x2, y2, z1), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x2, y2, z2), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(1.0, 0.0, 0.0)},
	}

	xm := NewMesh()
	xm.Vertices = []Vertex{
		{Position: lm.Vec3(x1, y1, z1), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x1, y1, z2), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x1, y2, z2), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x1, y1, z1), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x1, y2, z1), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
		{Position: lm.Vec3(x1, y2, z2), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(-1.0, 0.0, 0.0)},
	}

	zp := NewMesh()
	zp.Vertices = []Vertex{
		{Position: lm.Vec3(x1, y1, z2), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(x1, y2, z2), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(x2, y2, z2), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(x1, y1, z2), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(x2, y1, z2), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
		{Position: lm.Vec3(x2, y2, z2), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 0.0, 1.0)},
	}

	zm := NewMesh()
	zm.Vertices = []Vertex{
		{Position: lm.Vec3(x1, y1, z1), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(x1, y2, z1), Texcoord: lm.Vec2(0.0, 1.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(x2, y2, z1), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(x1, y1, z1), Texcoord: lm.Vec2(0.0, 0.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(x2, y1, z1), Texcoord: lm.Vec2(1.0, 0.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
		{Position: lm.Vec3(x2, y2, z1), Texcoord: lm.Vec2(1.0, 1.0), Normal: lm.Vec3(0.0, 0.0, -1.0)},
	}

	return []Mesh{yp, ym, xp, xm, zp, zm}
}

func Cube() *Model {
	model := NewModel()

	model.Meshes = append(model.Meshes, Cuboid(-0.5, -0.5, -0.5, 0.5, 0.5, 0.5)...)

	return &model
}
