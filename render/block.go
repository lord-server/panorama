package render

import (
	"image"

	"github.com/weqqr/panorama/game"
	"github.com/weqqr/panorama/world"
)

func overlayWithDepth(target *image.NRGBA, targetDepth *DepthBuffer, source *image.NRGBA, sourceDepth *DepthBuffer, origin image.Point) {
	width := source.Rect.Dx()
	height := source.Rect.Dy()

	for y := origin.Y; y < origin.Y+height; y++ {
		for x := origin.X; x < origin.X+width; x++ {
			c := source.NRGBAAt(x-origin.X, y-origin.Y)
			if c.A == 0 {
				// TODO: support opacity
				continue
			}
			target.SetNRGBA(x, y, c)
		}
	}
}

func RenderBlock(nr *NodeRasterizer, block *world.MapBlock, game *game.Game) *image.NRGBA {
	rect := image.Rect(0, 0, 16*BaseResolution, 16*BaseResolution+BaseResolution/8)
	img := image.NewNRGBA(rect)
	depth := NewDepthBuffer(rect)

	originX, originY := rect.Dx()/2-BaseResolution/2, rect.Dy()/2-BaseResolution/2-2

	for z := 0; z < world.MapBlockSize; z++ {
		for y := 0; y < world.MapBlockSize; y++ {
			for x := 0; x < world.MapBlockSize; x++ {
				node := block.GetNode(x, y, z)
				nodeName := block.ResolveName(node.ID)
				gameNode := game.Node(nodeName)

				nodeImg, nodeDepth := nr.Render(&gameNode)

				tileX, tileY := originX+BaseResolution*(z-x)/2, originY+BaseResolution*(z+x-2*y)/4

				overlayWithDepth(img, depth, nodeImg, nodeDepth, image.Pt(tileX, tileY))
			}
		}
	}

	return img
}
