package render

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/lm"
	"github.com/weqqr/panorama/mesh"
)

const BaseResolution = 64

type NodeRasterizer struct {
}

func NewNodeRasterizer() NodeRasterizer {
	return NodeRasterizer{}
}

func cartesianToBarycentric(p lm.Vector2, a, b, c lm.Vector2) lm.Vector3 {
	u := lm.Vec3(c.X-a.X, b.X-a.X, a.X-p.X)
	v := lm.Vec3(c.Y-a.Y, b.Y-a.Y, a.Y-p.Y)
	w := u.Cross(v)

	return lm.Vec3(1-(w.X+w.Y)/w.Z, w.Y/w.Z, w.X/w.Z)
}

func sampleTriangle(x, y int, a, b, c lm.Vector2) (bool, lm.Vector3) {
	p := lm.Vec2(float32(x), float32(y))

	samplePointOffset := lm.Vec2(0.5, 0.5)

	barycentric := cartesianToBarycentric(p.Add(samplePointOffset), a, b, c)

	if barycentric.X > 0 && barycentric.Y > 0 && barycentric.Z > 0 {
		return true, barycentric
	}

	return false, lm.Vector3{}
}

var LightDir lm.Vector3 = lm.Vec3(-0.9, 1, -0.7).Normalize()
var Projection lm.Matrix3 = lm.DimetricProjection()

func drawTriangle(img *image.NRGBA, depth *DepthBuffer, a, b, c mesh.Vertex) {
	originX := float32(img.Bounds().Dx() / 2)
	originY := float32(img.Bounds().Dy() / 2)
	origin := lm.Vec2(originX, originY)

	a.Position = Projection.MulVec(a.Position)
	b.Position = Projection.MulVec(b.Position)
	c.Position = Projection.MulVec(c.Position)

	pa := a.Position.XY().Mul(lm.Vec2(1, -1)).MulScalar(BaseResolution * math.Sqrt2 / 2).Add(origin)
	pb := b.Position.XY().Mul(lm.Vec2(1, -1)).MulScalar(BaseResolution * math.Sqrt2 / 2).Add(origin)
	pc := c.Position.XY().Mul(lm.Vec2(1, -1)).MulScalar(BaseResolution * math.Sqrt2 / 2).Add(origin)

	bboxMin := pa.Min(pb).Min(pc)
	bboxMax := pa.Max(pb).Max(pc)

	for y := int(bboxMin.Y); y < int(bboxMax.Y)+1; y++ {
		for x := int(bboxMin.X); x < int(bboxMax.X)+1; x++ {
			pointIsInsideTriangle, barycentric := sampleTriangle(x, y, pa, pb, pc)

			if !pointIsInsideTriangle {
				continue
			}

			pixelDepth := lm.Vec3(a.Position.Z, b.Position.Z, c.Position.Z).Dot(barycentric)

			if pixelDepth > depth.At(x, y) {
				continue
			}

			depth.Set(x, y, pixelDepth)

			normal := a.Normal.MulScalar(barycentric.X).
				Add(b.Normal.MulScalar(barycentric.Y)).
				Add(c.Normal.MulScalar(barycentric.Z))

			lighting := lm.Abs(lm.Clamp(normal.Dot(LightDir)*0.8+0.2, 0.0, 1.0))
			c := uint8(255 * lighting)

			finalColor := color.NRGBA{
				R: uint8(c),
				G: uint8(c),
				B: uint8(c),
				A: 255,
			}

			img.SetNRGBA(x, y, finalColor)
		}
	}
}

func (r *NodeRasterizer) Render(def *game.Node) (*image.NRGBA, *DepthBuffer) {
	rect := image.Rect(0, 0, BaseResolution, BaseResolution+BaseResolution/8-2)
	log.Printf("%v\n", rect)
	img := image.NewNRGBA(rect)
	depth := NewDepthBuffer(rect)

	if def.Mesh == nil {
		return img, depth
	}

	triangleCount := len(def.Mesh.Vertices) / 3

	for i := 0; i < triangleCount; i++ {
		a := def.Mesh.Vertices[i*3]
		b := def.Mesh.Vertices[i*3+1]
		c := def.Mesh.Vertices[i*3+2]

		drawTriangle(img, depth, a, b, c)
	}

	return img, depth
}
