package isometric

import (
	"github.com/weqqr/panorama/pkg/world"
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
