package mesh

import (
	"log"

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

func Cube() *Mesh {
	mesh, err := LoadOBJ("untitled.obj")
	if err != nil {
		log.Panic(err)
	}

	return &mesh
}
