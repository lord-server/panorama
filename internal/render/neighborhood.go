package render

import (
	"github.com/lord-server/panorama/internal/spatial"
	"github.com/lord-server/panorama/internal/world"
)

type BlockNeighborhood struct {
	blocks [27]*world.MapBlock
}

var neighborhoodCenter = spatial.BlockPosition{X: 1, Y: 1, Z: 1}

func blockIndex(pos spatial.BlockPosition) int {
	return pos.Z*9 + pos.Y*3 + pos.X
}

func (b *BlockNeighborhood) FetchBlock(w *world.World, posOffset, worldPos spatial.BlockPosition) {
	block, err := w.GetBlock(worldPos.Add(posOffset))

	if err != nil {
		return
	}

	b.SetBlock(neighborhoodCenter.Add(posOffset), block)
}

func (b *BlockNeighborhood) SetBlock(pos spatial.BlockPosition, block *world.MapBlock) {
	b.blocks[blockIndex(pos)] = block
}

func (b *BlockNeighborhood) getBlockByNodePos(pos spatial.NodePosition) *world.MapBlock {
	blockPos := spatial.BlockPosition{
		X: pos.X/spatial.BlockSize + neighborhoodCenter.X,
		Y: pos.Y/spatial.BlockSize + neighborhoodCenter.Y,
		Z: pos.Z/spatial.BlockSize + neighborhoodCenter.Z,
	}

	return b.blocks[blockIndex(blockPos)]
}

func (b *BlockNeighborhood) GetNode(pos spatial.NodePosition) (string, uint8, uint8) {
	block := b.getBlockByNodePos(pos)

	if block == nil {
		return "ignore", 0, 0
	}

	node := block.GetNode(spatial.NodePosition{
		X: pos.X % spatial.BlockSize,
		Y: pos.Y % spatial.BlockSize,
		Z: pos.Z % spatial.BlockSize,
	})
	name := block.ResolveName(node.ID)
	return name, node.Param1, node.Param2
}

func (b *BlockNeighborhood) GetParam1(pos spatial.NodePosition) uint8 {
	block := b.getBlockByNodePos(pos)

	if block == nil {
		return 0
	}

	node := block.GetNode(spatial.NodePosition{
		X: pos.X % spatial.BlockSize,
		Y: pos.Y % spatial.BlockSize,
		Z: pos.Z % spatial.BlockSize,
	})

	return node.Param1
}
