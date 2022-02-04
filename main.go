package main

import (
	"image"
	"image/png"
	"os"
)

const BaseResolution = 512

type NodeRasterizer struct {
}

func NewNodeRasterizer() NodeRasterizer {
	return NodeRasterizer{}
}

func (r *NodeRasterizer) RenderIsometric(def *NodeDef) *image.NRGBA {
	size := BaseResolution + BaseResolution/2
	canvas := image.NewNRGBA(image.Rect(0, 0, size, size))
	return canvas
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

func main() {
	nr := NewNodeRasterizer()
	mesh, err := loadOBJ("test.obj")
	if err != nil {
		panic(err)
	}
	def := &NodeDef{
		drawtype: DrawTypeMesh,
		mesh:     &mesh,
	}
	img := nr.RenderIsometric(def)
	savePNG(img, "test.png")
}
