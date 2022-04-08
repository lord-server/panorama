package isometric

import (
	"image"
	"math"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/raster"
	"github.com/weqqr/panorama/render"
	"github.com/weqqr/panorama/world"
)

const BaseResolution = 16

var (
	YOffsetCoef     = int(math.Round(BaseResolution * (1 + math.Sqrt2) / 4))
	TileBlockWidth  = world.MapBlockSize * BaseResolution
	TileBlockHeight = BaseResolution/2*world.MapBlockSize - 1 + YOffsetCoef*world.MapBlockSize
)

func decodeLight(param1 uint8) float32 {
	var LUT = [16]float32{
		0.000,
		0.024,
		0.059,
		0.118,
		0.196,
		0.286,
		0.384,
		0.471,
		0.545,
		0.608,
		0.659,
		0.710,
		0.769,
		0.835,
		0.918,
		1.000,
	}
	return LUT[param1&0xF]
}

type BlockNeighborhood struct {
	blocks [27]*world.MapBlock
}

func (b *BlockNeighborhood) FetchBlock(bx, by, bz, wx, wy, wz int, w *world.World) {
	block, err := w.GetBlock(wx, wy, wz)

	if err != nil {
		return
	}

	b.SetBlock(bx, by, bz, block)
}

func (b *BlockNeighborhood) SetBlock(bx, by, bz int, block *world.MapBlock) {
	b.blocks[bz*9+by*3+bx] = block
}

func (b *BlockNeighborhood) GetBlockAt(x, y, z int) *world.MapBlock {
	bx := x/16 + 1
	by := y/16 + 1
	bz := z/16 + 1

	return b.blocks[bz*9+by*3+bx]
}

func (b *BlockNeighborhood) GetNode(x, y, z int) (string, uint8, uint8) {
	block := b.GetBlockAt(x, y, z)

	if block == nil {
		return "air", 0, 0
	}

	node := block.GetNode(x%16, y%16, z%16)
	name := block.ResolveName(node.ID)
	return name, node.Param1, node.Param2
}

func (b *BlockNeighborhood) GetParam1(x, y, z int) uint8 {
	block := b.GetBlockAt(x, y, z)

	if block == nil {
		return 0
	}

	node := block.GetNode(x%16, y%16, z%16)
	return node.Param1
}

func renderBlock(target *raster.RenderBuffer, nr *render.NodeRasterizer, neighborhood *BlockNeighborhood, g *game.Game, offsetX, offsetY int, depth float32) {
	rect := image.Rect(0, 0, TileBlockWidth, TileBlockHeight)

	// FIXME: nodes must define their origin points
	// FIXME: Magic numbers and sloppy usage of BaseResolution
	originX, originY := rect.Dx()/2-BaseResolution/2, rect.Dy()/2+BaseResolution/4+2

	for z := world.MapBlockSize - 1; z >= 0; z-- {
		for y := world.MapBlockSize - 1; y >= 0; y-- {
			for x := world.MapBlockSize - 1; x >= 0; x-- {
				tileOffsetX := originX + BaseResolution*(z-x)/2 + offsetX
				tileOffsetY := originY + BaseResolution/4*(z+x) - YOffsetCoef*y + offsetY

				// Fast path: Don't bother with nodes outside viewport
				nodeTileTooLow := tileOffsetX <= target.Color.Rect.Min.X-BaseResolution || tileOffsetY <= target.Color.Rect.Min.Y-BaseResolution-BaseResolution/8
				nodeTileTooHigh := tileOffsetX >= target.Color.Rect.Max.X || tileOffsetY >= target.Color.Rect.Max.Y

				if nodeTileTooLow || nodeTileTooHigh {
					continue
				}

				name, param1, param2 := neighborhood.GetNode(x, y, z)

				// Fast path: checking for air immediately is faster than fetching NodeDefinition
				if name == "air" {
					continue
				}

				nodeDef := g.NodeDef(name)

				light := decodeLight(param1)
				if l := decodeLight(neighborhood.GetParam1(x+1, y, z)); l > light {
					light = l
				}
				if l := decodeLight(neighborhood.GetParam1(x, y+1, z)); l > light {
					light = l
				}
				if l := decodeLight(neighborhood.GetParam1(x, y, z+1)); l > light {
					light = l
				}

				renderableNode := render.RenderableNode{
					Name:   name,
					Light:  light,
					Param2: param2,
				}
				renderedNode := nr.Render(renderableNode, &nodeDef)

				depthOffset := -float32(z+x)/math.Sqrt2 - 0.5*(float32(y)) + depth
				target.OverlayDepthAware(renderedNode, image.Pt(tileOffsetX, tileOffsetY), depthOffset)
			}
		}
	}
}
