package imaging

import "image"

func OverlayWithDepth(target *image.NRGBA, targetDepth *Depth, source *image.NRGBA, sourceDepth *Depth, origin image.Point, depthOffset float32) {
	if source == nil {
		return
	}

	bbox := source.Rect.Add(origin).Intersect(target.Rect)

	for y := bbox.Min.Y; y < bbox.Max.Y; y++ {
		for x := bbox.Min.X; x < bbox.Max.X; x++ {
			targetZ := targetDepth.At(x, y)
			sourceZ := sourceDepth.At(x-origin.X, y-origin.Y) + depthOffset

			if sourceZ > targetZ {
				continue
			}

			targetDepth.Set(x, y, sourceZ)

			c := source.NRGBAAt(x-origin.X, y-origin.Y)
			if c.A == 0 {
				// TODO: support opacity
				continue
			}
			target.SetNRGBA(x, y, c)
		}
	}
}
