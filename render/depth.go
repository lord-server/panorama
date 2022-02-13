package render

import (
	"image"
	"image/png"
	"math"
	"os"
)

type DepthBuffer struct {
	Pix  []float32
	Rect image.Rectangle
}

func NewDepthBuffer(rect image.Rectangle) *DepthBuffer {
	pix := make([]float32, rect.Dx()*rect.Dy())
	for i := range pix {
		pix[i] = math.MaxFloat32
	}

	return &DepthBuffer{
		Pix:  pix,
		Rect: rect,
	}
}

func (d *DepthBuffer) At(x, y int) float32 {
	if x < d.Rect.Min.X || y < d.Rect.Min.Y || x >= d.Rect.Max.X || y >= d.Rect.Max.Y {
		return -math.MaxFloat32
	}
	return d.Pix[d.Rect.Dx()*y+x]
}

func (d *DepthBuffer) Set(x, y int, depth float32) {
	if x < d.Rect.Min.X || y < d.Rect.Min.Y || x > d.Rect.Max.X || y > d.Rect.Max.Y {
		return
	}
	d.Pix[d.Rect.Dx()*y+x] = depth
}

func savePNG(img *image.NRGBA, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}

	if err := png.Encode(file, img); err != nil {
		file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
