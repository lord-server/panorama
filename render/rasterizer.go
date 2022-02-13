package render

import (
	"image"
	"image/color"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/lm"
	"github.com/weqqr/panorama/mesh"
)

const BaseResolution = 64

var (
	YOffsetCoef     = int(math.Round(BaseResolution * (1 + math.Sqrt2) / 4))
	TileBlockWidth  = 16 * BaseResolution
	TileBlockHeight = BaseResolution/2*16 - 1 + YOffsetCoef*16
)

type NodeRasterizer struct {
	cache      map[string]*image.NRGBA
	depthCache map[string]*DepthBuffer
}

func NewNodeRasterizer() NodeRasterizer {
	return NodeRasterizer{
		cache:      make(map[string]*image.NRGBA),
		depthCache: make(map[string]*DepthBuffer),
	}
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

func sampleTexture(tex *image.NRGBA, texcoord lm.Vector2) lm.Vector3 {
	x := int(texcoord.X * float32(tex.Rect.Dx()))
	y := int(texcoord.Y * float32(tex.Rect.Dy()))
	c := tex.NRGBAAt(x, y)
	return lm.Vec3(float32(c.R)/255, float32(c.G)/255, float32(c.B)/255)
}

var LightDir lm.Vector3 = lm.Vec3(-0.9, 1, -0.7).Normalize()
var Projection lm.Matrix3 = lm.DimetricProjection()

func drawTriangle(img *image.NRGBA, depth *DepthBuffer, tex *image.NRGBA, a, b, c mesh.Vertex) {
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

			var finalColor color.NRGBA
			if tex != nil {
				texcoord := a.Texcoord.MulScalar(barycentric.X).
					Add(b.Texcoord.MulScalar(barycentric.Y)).
					Add(c.Texcoord.MulScalar(barycentric.Z))
				col := sampleTexture(tex, texcoord).MulScalar(lighting)

				finalColor = color.NRGBA{
					R: uint8(255 * col.X),
					G: uint8(255 * col.Y),
					B: uint8(255 * col.Z),
					A: 255,
				}
			} else {
				finalColor = color.NRGBA{
					R: uint8(255 * lighting),
					G: uint8(255 * lighting),
					B: uint8(255 * lighting),
					A: 255,
				}
			}

			img.SetNRGBA(x, y, finalColor)
		}
	}
}

func (r *NodeRasterizer) Render(wnode string, node *game.Node) (*image.NRGBA, *DepthBuffer) {
	if img, ok := r.cache[wnode]; ok {
		return img, r.depthCache[wnode]
	}

	rect := image.Rect(0, 0, BaseResolution, BaseResolution+BaseResolution/8-2)
	img := image.NewNRGBA(rect)
	depth := NewDepthBuffer(rect)
	if node.DrawType == game.DrawTypeAirLlke {
		return img, depth
	}

	if node.Mesh == nil {
		return img, depth
	}

	triangleCount := len(node.Mesh.Vertices) / 3

	var texture *image.NRGBA
	if len(node.Tiles) >= 1 {
		texture = node.Tiles[0]
	}

	for i := 0; i < triangleCount; i++ {
		a := node.Mesh.Vertices[i*3]
		b := node.Mesh.Vertices[i*3+1]
		c := node.Mesh.Vertices[i*3+2]

		drawTriangle(img, depth, texture, a, b, c)
	}

	r.cache[wnode] = img
	r.depthCache[wnode] = depth

	return img, depth
}
