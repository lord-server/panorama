package game

import (
	"image"
	"image/color"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/weqqr/panorama/mesh"
)

type MediaCache struct {
	images     map[string]*image.NRGBA
	meshes     map[string]*mesh.Mesh
	dummyImage *image.NRGBA
}

func NewMediaCache() *MediaCache {
	dummyImage := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	dummyImage.SetNRGBA(0, 0, color.NRGBA{255, 0, 255, 255})
	dummyImage.SetNRGBA(0, 1, color.NRGBA{0, 0, 0, 255})
	dummyImage.SetNRGBA(1, 0, color.NRGBA{0, 0, 0, 255})
	dummyImage.SetNRGBA(1, 1, color.NRGBA{255, 0, 255, 255})

	return &MediaCache{
		images:     make(map[string]*image.NRGBA),
		meshes:     make(map[string]*mesh.Mesh),
		dummyImage: dummyImage,
	}
}

func (m *MediaCache) fetchMedia(path string) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsRegular() {
			return nil
		}

		basePath := filepath.Base(path)
		switch filepath.Ext(path) {
		case ".png":
			img, _ := LoadPNG(path)
			m.images[basePath] = img
		case ".obj":
			log.Println(path)
			mesh, err := mesh.LoadOBJ(path)
			if err != nil {
				return err
			}
			m.meshes[basePath] = &mesh
		}

		return nil
	})
}

func (m *MediaCache) Image(name string) *image.NRGBA {
	if img, ok := m.images[name]; ok {
		return img
	} else {
		log.Printf("unknown image: %v\n", name)
		return m.dummyImage
	}
}

func (m *MediaCache) Mesh(name string) *mesh.Mesh {
	if mesh, ok := m.meshes[name]; ok {
		return mesh
	} else {
		log.Printf("unknown mesh: %v\n", name)
		return nil
	}
}
