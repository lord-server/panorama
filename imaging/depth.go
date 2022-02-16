package imaging

import (
	"image"
	"math"
)

type Depth struct {
	Pix  []float32
	Rect image.Rectangle
}

func NewDepth(rect image.Rectangle) *Depth {
	pix := make([]float32, rect.Dx()*rect.Dy())
	for i := range pix {
		pix[i] = math.MaxFloat32
	}

	return &Depth{
		Pix:  pix,
		Rect: rect,
	}
}

func (d *Depth) At(x, y int) float32 {
	if x < d.Rect.Min.X || y < d.Rect.Min.Y || x >= d.Rect.Max.X || y >= d.Rect.Max.Y {
		return -math.MaxFloat32
	}
	return d.Pix[d.Rect.Dx()*y+x]
}

func (d *Depth) Set(x, y int, depth float32) {
	if x < d.Rect.Min.X || y < d.Rect.Min.Y || x > d.Rect.Max.X || y > d.Rect.Max.Y {
		return
	}
	d.Pix[d.Rect.Dx()*y+x] = depth
}
