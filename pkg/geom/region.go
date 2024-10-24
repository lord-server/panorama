package geom

// Bounds defines the extent of a region on a single axis. It is assumed
// that Max is always greater or equal to Min.
type Bounds struct {
	Min int `toml:"min"`
	Max int `toml:"max"`
}

// Region defines an axis-aligned cuboid region in world space (units are nodes).
type Region struct {
	XBounds Bounds `toml:"x_bounds"`
	YBounds Bounds `toml:"y_bounds"`
	ZBounds Bounds `toml:"z_bounds"`
}

func (lhs Region) Intersects(rhs Region) bool {
	xOverlaps := lhs.XBounds.Min <= rhs.XBounds.Max && rhs.XBounds.Min <= lhs.XBounds.Max
	yOverlaps := lhs.YBounds.Min <= rhs.YBounds.Max && rhs.YBounds.Min <= lhs.YBounds.Max
	zOverlaps := lhs.ZBounds.Min <= rhs.ZBounds.Max && rhs.ZBounds.Min <= lhs.ZBounds.Max

	return xOverlaps && yOverlaps && zOverlaps
}

func (lhs Region) IsAtEdge(pos NodePosition) bool {
	isAtXEdge := pos.X == lhs.XBounds.Max || pos.X == lhs.XBounds.Min
	isAtYEdge := pos.Y == lhs.YBounds.Max || pos.Y == lhs.YBounds.Min
	isAtZEdge := pos.Z == lhs.ZBounds.Max || pos.Z == lhs.ZBounds.Min

	return isAtXEdge || isAtYEdge || isAtZEdge
}

// ProjectedRegion defines an axis-aligned rectangle region in tile space (units are
// tiles at zoom level 0). It's used to represent a projection of a Region onto
// the screen.
type ProjectedRegion struct {
	XBounds Bounds
	YBounds Bounds
}
