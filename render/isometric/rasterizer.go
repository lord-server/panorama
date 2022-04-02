package isometric

import (
	"image"
	"image/color"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/lm"
	"github.com/weqqr/panorama/mesh"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/world"
)

const BaseResolution = 16

var (
	YOffsetCoef     = int(math.Round(BaseResolution * (1 + math.Sqrt2) / 4))
	TileBlockWidth  = world.MapBlockSize * BaseResolution
	TileBlockHeight = BaseResolution/2*world.MapBlockSize - 1 + YOffsetCoef*world.MapBlockSize
)

type RenderableNode struct {
	Name   string
	Light  float32
	Param2 uint8
}

type NodeRasterizer struct {
	cache map[RenderableNode]*raster.RenderBuffer
}

func NewNodeRasterizer() NodeRasterizer {
	return NodeRasterizer{
		cache: make(map[RenderableNode]*raster.RenderBuffer),
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

func sampleTexture(tex *image.NRGBA, texcoord lm.Vector2) lm.Vector4 {
	x := int(texcoord.X * float32(tex.Rect.Dx()))
	y := int(texcoord.Y * float32(tex.Rect.Dy()))
	c := tex.NRGBAAt(x, y)
	return lm.Vec4(float32(c.R)/255, float32(c.G)/255, float32(c.B)/255, float32(c.A)/255)
}

var LightDir = lm.Vec3(-0.6, 1, -0.8).Normalize()
var LightIntensity = 0.95 / LightDir.MaxComponent()
var Projection = lm.DimetricProjection()

func drawTriangle(target *raster.RenderBuffer, tex *image.NRGBA, light float32, a, b, c mesh.Vertex) {
	originX := float32(target.Color.Bounds().Dx() / 2)
	originY := float32(target.Color.Bounds().Dy() / 2)
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

			normal := a.Normal.MulScalar(barycentric.X).
				Add(b.Normal.MulScalar(barycentric.Y)).
				Add(c.Normal.MulScalar(barycentric.Z))

			lighting := LightIntensity * light * lm.Clamp(lm.Abs(normal.Dot(LightDir))*0.8+0.2, 0.0, 1.0)

			var finalColor color.NRGBA
			if tex != nil {
				texcoord := a.Texcoord.MulScalar(barycentric.X).
					Add(b.Texcoord.MulScalar(barycentric.Y)).
					Add(c.Texcoord.MulScalar(barycentric.Z))
				rgba := sampleTexture(tex, texcoord)
				col := rgba.XYZ().MulScalar(lighting).ClampScalar(0.0, 1.0)

				finalColor = color.NRGBA{
					R: uint8(255 * col.X),
					G: uint8(255 * col.Y),
					B: uint8(255 * col.Z),
					A: uint8(255 * rgba.W),
				}
			} else {
				finalColor = color.NRGBA{
					R: uint8(255 * lighting),
					G: uint8(255 * lighting),
					B: uint8(255 * lighting),
					A: 255,
				}
			}

			if finalColor.A > 10 {
				if pixelDepth > target.Depth.At(x, y) {
					continue
				}

				target.Color.SetNRGBA(x, y, finalColor)
				target.Depth.Set(x, y, pixelDepth)
			}
		}
	}
}

func (r *NodeRasterizer) Render(node RenderableNode, nodeDef *game.NodeDefinition) *raster.RenderBuffer {
	if nodeDef.DrawType == game.DrawTypeAirlike || nodeDef.Model == nil || len(nodeDef.Textures) == 0 {
		return nil
	}

	if target, ok := r.cache[node]; ok {
		return target
	}

	rect := image.Rect(0, 0, BaseResolution, BaseResolution+BaseResolution/8)
	target := raster.NewRenderBuffer(rect)

	for j, mesh := range nodeDef.Model.Meshes {
		triangleCount := len(mesh.Vertices) / 3

		for i := 0; i < triangleCount; i++ {
			a := mesh.Vertices[i*3]
			b := mesh.Vertices[i*3+1]
			c := mesh.Vertices[i*3+2]

			drawTriangle(target, nodeDef.Textures[j], node.Light, a, b, c)
		}
	}

	r.cache[node] = target

	return target
}
