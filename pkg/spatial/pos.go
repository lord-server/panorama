package spatial

// NodePos is a node position in world space
type NodePos struct {
	X, Y, Z int
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

func (lhs BlockPos) Add(rhs BlockPos) BlockPos {
	return BlockPos{
		X: lhs.X + rhs.X,
		Y: lhs.Y + rhs.Y,
		Z: lhs.Z + rhs.Z,
	}
}
