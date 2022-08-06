package spatial

const BlockSize = 16
const BlockVolume = BlockSize * BlockSize * BlockSize

// NodePos is a node position in world space
type NodePos struct {
	X, Y, Z int
}

func (lhs NodePos) Region() Region {
	return Region{
		XBounds: Bounds{Min: lhs.X, Max: lhs.X},
		YBounds: Bounds{Min: lhs.Y, Max: lhs.Y},
		ZBounds: Bounds{Min: lhs.Z, Max: lhs.Z},
	}
}

func (lhs NodePos) Add(rhs NodePos) NodePos {
	return NodePos{
		X: lhs.X + rhs.X,
		Y: lhs.Y + rhs.Y,
		Z: lhs.Z + rhs.Z,
	}
}

// BlockPos is a block position in world space
type BlockPos struct {
	X, Y, Z int
}

func (lhs BlockPos) AddNode(pos NodePos) NodePos {
	return NodePos{
		X: lhs.X*BlockSize + pos.X,
		Y: lhs.Y*BlockSize + pos.Y,
		Z: lhs.Z*BlockSize + pos.Z,
	}
}

func (lhs BlockPos) Add(rhs BlockPos) BlockPos {
	return BlockPos{
		X: lhs.X + rhs.X,
		Y: lhs.Y + rhs.Y,
		Z: lhs.Z + rhs.Z,
	}
}
