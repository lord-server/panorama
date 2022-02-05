package main

import (
	"image"
	"image/color"
	"math"
)

const BaseResolution = 64

type NodeRasterizer struct {
}

func NewNodeRasterizer() NodeRasterizer {
	return NodeRasterizer{}
}

func cartesianToBarycentric(p Vector2, a, b, c Vector2) Vector3 {
	u := NewVector3(c.X-a.X, b.X-a.X, a.X-p.X)
	v := NewVector3(c.Y-a.Y, b.Y-a.Y, a.Y-p.Y)
	w := u.Cross(v)

	return NewVector3(1-(w.X+w.Y)/w.Z, w.Y/w.Z, w.X/w.Z)
}

func sampleTriangle(x, y int, a, b, c Vector2) (bool, Vector3) {
	p := NewVector2(float32(x), float32(y))

	samplePointOffsets := []Vector2{
		{0.5, 0.5},
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}

	for _, offset := range samplePointOffsets {
		barycentric := cartesianToBarycentric(p.Add(offset), a, b, c)

		if barycentric.X > 0 && barycentric.Y > 0 && barycentric.Z > 0 {
			return true, barycentric
		}
	}

	return false, Vector3{}
}

var Projection Matrix3 = DimetricProjection()

func drawTriangle(img *image.NRGBA, a, b, c Vertex) {
	originX := float32(img.Bounds().Dx() / 2)
	originY := float32(img.Bounds().Dy() / 2)
	origin := NewVector2(originX, originY)

	pa := Projection.MulVec(a.position).XY().MulScalar(BaseResolution * math.Sqrt2 / 2).Add(origin)
	pb := Projection.MulVec(b.position).XY().MulScalar(BaseResolution * math.Sqrt2 / 2).Add(origin)
	pc := Projection.MulVec(c.position).XY().MulScalar(BaseResolution * math.Sqrt2 / 2).Add(origin)

	bboxMin := pa.Min(pb).Min(pc)
	bboxMax := pa.Max(pb).Max(pc)

	for y := int(bboxMin.Y); y < int(bboxMax.Y)+1; y++ {
		for x := int(bboxMin.X); x < int(bboxMax.X)+1; x++ {
			pointIsInsideTriangle, _ := sampleTriangle(x, y, pa, pb, pc)

			if !pointIsInsideTriangle {
				continue
			}

			img.SetNRGBA(x, y, color.NRGBA{255, 0, 0, 255})
		}
	}
}

func (r *NodeRasterizer) Render(def *NodeDef) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, BaseResolution, BaseResolution+BaseResolution/8))

	triangleCount := len(def.mesh.vertices) / 3

	for i := 0; i < triangleCount; i++ {
		a := def.mesh.vertices[i*3]
		b := def.mesh.vertices[i*3+1]
		c := def.mesh.vertices[i*3+2]
		drawTriangle(img, a, b, c)
	}

	return img
}
