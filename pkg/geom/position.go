package geom

const BlockSize = 16
const BlockVolume = BlockSize * BlockSize * BlockSize

// NodePosition is a node position in world space
type NodePosition struct {
	X, Y, Z int
}

func (lhs NodePosition) Region() Region {
	return Region{
		XBounds: Bounds{Min: lhs.X, Max: lhs.X},
		YBounds: Bounds{Min: lhs.Y, Max: lhs.Y},
		ZBounds: Bounds{Min: lhs.Z, Max: lhs.Z},
	}
}

func (lhs NodePosition) Add(rhs NodePosition) NodePosition {
	return NodePosition{
		X: lhs.X + rhs.X,
		Y: lhs.Y + rhs.Y,
		Z: lhs.Z + rhs.Z,
	}
}

// BlockPosition is a block position in world space
type BlockPosition struct {
	X, Y, Z int
}

func (lhs BlockPosition) AddNode(pos NodePosition) NodePosition {
	return NodePosition{
		X: lhs.X*BlockSize + pos.X,
		Y: lhs.Y*BlockSize + pos.Y,
		Z: lhs.Z*BlockSize + pos.Z,
	}
}

func (lhs BlockPosition) Add(rhs BlockPosition) BlockPosition {
	return BlockPosition{
		X: lhs.X + rhs.X,
		Y: lhs.Y + rhs.Y,
		Z: lhs.Z + rhs.Z,
	}
}
