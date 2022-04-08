package render

import (
	"image"
	"image/color"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/lm"
	"github.com/weqqr/panorama/mesh"
	"github.com/weqqr/panorama/raster"
)

const Gamma = 2.2

type RenderableNode struct {
	Name   string
	Light  float32
	Param2 uint8
}

type NodeRasterizer struct {
	targetWidth  int
	targetHeight int
	scale        float32
	projection   lm.Matrix3

	cache map[RenderableNode]*raster.RenderBuffer
}

func NewNodeRasterizer(width, height int, scale float32, projection lm.Matrix3) NodeRasterizer {
	return NodeRasterizer{
		targetWidth:  width,
		targetHeight: height,
		scale:        scale,
		projection:   projection,

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

var SunLightDir = lm.Vec3(-0.5, 1, -0.8).Normalize()
var SunLightIntensity = 0.95 / SunLightDir.MaxComponent()

func (r *NodeRasterizer) drawTriangle(target *raster.RenderBuffer, tex *image.NRGBA, lighting float32, a, b, c mesh.Vertex) {
	originX := float32(target.Color.Bounds().Dx() / 2)
	originY := float32(target.Color.Bounds().Dy() / 2)
	origin := lm.Vec2(originX, originY)

	a.Position = r.projection.MulVec(a.Position)
	b.Position = r.projection.MulVec(b.Position)
	c.Position = r.projection.MulVec(c.Position)

	pa := a.Position.XY().Mul(lm.Vec2(1, -1)).MulScalar(r.scale).Add(origin)
	pb := b.Position.XY().Mul(lm.Vec2(1, -1)).MulScalar(r.scale).Add(origin)
	pc := c.Position.XY().Mul(lm.Vec2(1, -1)).MulScalar(r.scale).Add(origin)

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

			lighting := SunLightIntensity * lighting * lm.Clamp(lm.Abs(normal.Dot(SunLightDir))*0.8+0.2, 0.0, 1.0)

			var finalColor color.NRGBA
			if tex != nil {
				texcoord := a.Texcoord.MulScalar(barycentric.X).
					Add(b.Texcoord.MulScalar(barycentric.Y)).
					Add(c.Texcoord.MulScalar(barycentric.Z))
				rgba := sampleTexture(tex, texcoord)
				col := rgba.XYZ().PowScalar(Gamma).MulScalar(lighting).PowScalar(1.0/Gamma).ClampScalar(0.0, 1.0)

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

func transformToFaceDir(v lm.Vector3, facedir uint8) lm.Vector3 {
	axis := (facedir >> 2) & 0x7
	dir := facedir & 0x3

	// Left click with screwdriver
	switch dir {
	case 0: // no-op
	case 1:
		v = v.RotateXZ(lm.Radians(-90))
	case 2:
		v = v.RotateXZ(lm.Radians(180))
	case 3:
		v = v.RotateXZ(lm.Radians(90))
	}

	// Right click with screwdriver
	switch axis {
	case 0: // no-op
	case 1:
		v = v.RotateYZ(lm.Radians(90))
	case 2:
		v = v.RotateYZ(lm.Radians(-90))
	case 3:
		v = v.RotateXY(lm.Radians(-90))
	case 4:
		v = v.RotateXY(lm.Radians(90))
	case 5:
		v = v.RotateXY(lm.Radians(180))
	}

	return v
}

func (r *NodeRasterizer) Render(node RenderableNode, nodeDef *game.NodeDefinition) *raster.RenderBuffer {
	if nodeDef.DrawType == game.DrawTypeAirlike || nodeDef.Model == nil || len(nodeDef.Textures) == 0 {
		return nil
	}

	if target, ok := r.cache[node]; ok {
		return target
	}

	rect := image.Rect(0, 0, r.targetWidth, r.targetHeight)
	target := raster.NewRenderBuffer(rect)

	for j, mesh := range nodeDef.Model.Meshes {
		triangleCount := len(mesh.Vertices) / 3

		for i := 0; i < triangleCount; i++ {
			a := mesh.Vertices[i*3]
			b := mesh.Vertices[i*3+1]
			c := mesh.Vertices[i*3+2]

			if nodeDef.ParamType2 == game.ParamType2FaceDir {
				a.Position = transformToFaceDir(a.Position, node.Param2)
				b.Position = transformToFaceDir(b.Position, node.Param2)
				c.Position = transformToFaceDir(c.Position, node.Param2)
				a.Normal = transformToFaceDir(a.Normal, node.Param2)
				b.Normal = transformToFaceDir(b.Normal, node.Param2)
				c.Normal = transformToFaceDir(c.Normal, node.Param2)
			}

			a.Position.Z = -a.Position.Z
			b.Position.Z = -b.Position.Z
			c.Position.Z = -c.Position.Z

			a.Position.X = -a.Position.X
			b.Position.X = -b.Position.X
			c.Position.X = -c.Position.X

			r.drawTriangle(target, nodeDef.Textures[j], node.Light, a, b, c)
		}
	}

	r.cache[node] = target

	return target
}
