package raster

import (
	"image"
	"image/color"
)

type RenderBuffer struct {
	Color *image.NRGBA
	Depth *Depth
	Dirty bool
}

func NewRenderBuffer(rect image.Rectangle) *RenderBuffer {
	return &RenderBuffer{
		Color: image.NewNRGBA(rect),
		Depth: NewDepth(rect),
		Dirty: false,
	}
}

func (target *RenderBuffer) OverlayDepthAwareWithAlpha(source *RenderBuffer, origin image.Point, depthOffset float64) {
	target.Dirty = true
	if source == nil {
		return
	}

	bbox := source.Color.Rect.Add(origin).Intersect(target.Color.Rect)

	for y := bbox.Min.Y; y < bbox.Max.Y; y++ {
		for x := bbox.Min.X; x < bbox.Max.X; x++ {
			targetZ := target.Depth.At(x, y)
			sourceZ := source.Depth.At(x-origin.X, y-origin.Y) + depthOffset

			if sourceZ > targetZ {
				continue
			}

			target.Depth.Set(x, y, sourceZ)

			c := source.Color.NRGBAAt(x-origin.X, y-origin.Y)
			if c.A == 0 {
				continue
			}

			d := target.Color.NRGBAAt(x, y)

			// Blend with alpha
			sourceA := float64(c.A) / 255
			targetA := float64(d.A) / 255

			outA := sourceA + targetA*(1-sourceA)
			outR := (float64(c.R)*sourceA + float64(d.R)*targetA*(1-sourceA)) / outA
			outG := (float64(c.G)*sourceA + float64(d.G)*targetA*(1-sourceA)) / outA
			outB := (float64(c.B)*sourceA + float64(d.B)*targetA*(1-sourceA)) / outA

			target.Color.SetNRGBA(x, y, color.NRGBA{
				R: uint8(outR),
				G: uint8(outG),
				B: uint8(outB),
				A: uint8(outA * 255),
			})
		}
	}
}

func (target *RenderBuffer) OverlayDepthAware(source *RenderBuffer, origin image.Point, depthOffset float64) {
	target.Dirty = true

	if source == nil {
		return
	}

	bbox := source.Color.Rect.Add(origin).Intersect(target.Color.Rect)

	// This loop is by far the hottest in the entire program.
	// All function calls and pixel offset calculations are
	// inlined and re-used to improve performance. Writing it
	// this way makes rendering about 20% faster compared to
	// naive implementation.
	for y := bbox.Min.Y; y < bbox.Max.Y; y++ {
		sourcePixelBaseOffset := (y-origin.Y)*source.Depth.Rect.Max.X - origin.X
		targetPixelBaseOffset := y * target.Depth.Rect.Max.X
		for x := bbox.Min.X; x < bbox.Max.X; x++ {
			sourcePixelOffset := sourcePixelBaseOffset + x
			targetPixelOffset := targetPixelBaseOffset + x

			sourceZ := source.Depth.Pix[sourcePixelOffset] + depthOffset
			targetZ := target.Depth.Pix[targetPixelOffset]

			if sourceZ > targetZ {
				continue
			}

			target.Depth.Pix[targetPixelOffset] = sourceZ

			sourcePixelOffset *= 4
			if source.Color.Pix[sourcePixelOffset+3] == 0 {
				// TODO: support opacity
				continue
			}

			targetPixelOffset *= 4

			target.Color.Pix[targetPixelOffset+0] = source.Color.Pix[sourcePixelOffset+0]
			target.Color.Pix[targetPixelOffset+1] = source.Color.Pix[sourcePixelOffset+1]
			target.Color.Pix[targetPixelOffset+2] = source.Color.Pix[sourcePixelOffset+2]
			target.Color.Pix[targetPixelOffset+3] = 255
		}
	}
}
