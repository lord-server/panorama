package rasterizer

import (
	"image"
	"math"
)

type Depth struct {
	Pix  []float64
	Rect image.Rectangle
}

func NewDepth(rect image.Rectangle) *Depth {
	pix := make([]float64, rect.Dx()*rect.Dy())

	for i := range pix {
		pix[i] = math.MaxFloat64
	}

	return &Depth{
		Pix:  pix,
		Rect: rect,
	}
}

func (d *Depth) At(x, y int) float64 {
	if x < d.Rect.Min.X || y < d.Rect.Min.Y || x >= d.Rect.Max.X || y >= d.Rect.Max.Y {
		return -math.MaxFloat64
	}

	return d.Pix[d.Rect.Dx()*y+x]
}

func (d *Depth) Set(x, y int, depth float64) {
	if x < d.Rect.Min.X || y < d.Rect.Min.Y || x > d.Rect.Max.X || y > d.Rect.Max.Y {
		return
	}

	d.Pix[d.Rect.Dx()*y+x] = depth
}
