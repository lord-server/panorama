package imaging

import "image"

type RenderBuffer struct {
	Color *image.NRGBA
	Depth *Depth
}

func NewRenderBuffer(rect image.Rectangle) *RenderBuffer {
	return &RenderBuffer{
		Color: image.NewNRGBA(rect),
		Depth: NewDepth(rect),
	}
}

func (target *RenderBuffer) OverlayDepthAware(source *RenderBuffer, origin image.Point, depthOffset float32) {
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
				// TODO: support opacity
				continue
			}
			target.Color.SetNRGBA(x, y, c)
		}
	}
}
