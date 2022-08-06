package region

// Bounds defines the extent of a region on a single axis. It is assumed
// that Max is always greater or equal to Min.
type Bounds struct {
	Min int `toml:"min"`
	Max int `toml:"max"`
}

// Region defines a cuboid region in world space (units are nodes).
type Region struct {
	XBounds Bounds `toml:"x_bounds"`
	YBounds Bounds `toml:"y_bounds"`
	ZBounds Bounds `toml:"z_bounds"`
}

// TileRegion defines a rectangle region in tile space (units are tiles at
// zoom level 0). It's used to represent a projection of a Region onto the
// screen.
type TileRegion struct {
	XBounds Bounds
	YBounds Bounds
}
