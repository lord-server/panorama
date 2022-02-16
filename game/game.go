package game

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/weqqr/panorama/mesh"
)

type GameDescriptor struct {
	Aliases map[string]string         `json:"aliases"`
	Nodes   map[string]NodeDescriptor `json:"nodes"`
}

func toNRGBA(img image.Image) *image.NRGBA {
	dst := image.NewNRGBA(img.Bounds())
	draw.Draw(dst, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return dst
}

func LoadPNG(path string) (*image.NRGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return toNRGBA(img), nil
}

type Node struct {
	DrawType DrawType
	Textures []*image.NRGBA
	Model    *mesh.Model
}

type Game struct {
	Aliases map[string]string
	Nodes   map[string]Node
	unknown Node
}

func makeNormalNode(drawtype DrawType, tiles []*image.NRGBA) Node {
	textures := make([]*image.NRGBA, 6)
	model := mesh.Cube()

	if len(tiles) == 0 {
		return Node{
			DrawType: drawtype,
			Textures: textures,
			Model:    model,
		}
	}

	for i := 0; i < 6; i++ {
		if i >= len(tiles) {
			textures[i] = tiles[len(tiles)-1]
			continue
		}

		textures[i] = tiles[i]
	}

	return Node{
		DrawType: drawtype,
		Textures: textures,
		Model:    model,
	}
}

func makeMeshNode(drawtype DrawType, model *mesh.Model, tiles []*image.NRGBA) Node {
	textures := make([]*image.NRGBA, len(model.Meshes))
	if len(tiles) == 0 {
		return Node{
			DrawType: drawtype,
			Textures: textures,
			Model:    model,
		}
	}

	for i := range model.Meshes {
		if i >= len(tiles) {
			break
		}
		textures[i] = tiles[i]
	}

	return Node{
		DrawType: drawtype,
		Model:    model,
		Textures: textures,
	}
}

func ResolveNode(descriptor NodeDescriptor, mediaCache *MediaCache, name string) Node {
	tiles := make([]*image.NRGBA, len(descriptor.Tiles))

	for i, tileName := range descriptor.Tiles {
		tiles[i] = mediaCache.Image(tileName)
	}

	switch descriptor.DrawType {
	case DrawTypeNormal, DrawTypeAllFaces, DrawTypeLiquid, DrawTypeFlowingLiquid, DrawTypeGlasslike, DrawTypeGlasslikeFramed:
		return makeNormalNode(descriptor.DrawType, tiles)
	case DrawTypeMesh:
		if descriptor.Mesh == nil {
			break
		}

		model := mediaCache.Mesh(*descriptor.Mesh)
		return makeMeshNode(descriptor.DrawType, model, tiles)
	}

	return Node{
		DrawType: descriptor.DrawType,
		Textures: []*image.NRGBA{},
		Model:    nil,
	}
}

func LoadGame(desc string, path string) (Game, error) {
	descJSON, err := os.ReadFile(desc)
	if err != nil {
		return Game{}, err
	}

	var descriptor GameDescriptor
	err = json.Unmarshal(descJSON, &descriptor)
	if err != nil {
		return Game{}, err
	}

	mediaCache := NewMediaCache()

	err = mediaCache.fetchMedia(path)
	if err != nil {
		return Game{}, err
	}

	nodes := make(map[string]Node)
	for name, gameNode := range descriptor.Nodes {
		node := ResolveNode(gameNode, mediaCache, name)

		nodes[name] = node
	}

	return Game{
		Aliases: descriptor.Aliases,
		Nodes:   nodes,
		unknown: Node{
			DrawType: DrawTypeNormal,
			Textures: []*image.NRGBA{mediaCache.dummyImage},
			Model:    nil,
		},
	}, nil
}

func (g *Game) Node(node string) Node {
	if def, ok := g.Nodes[node]; ok {
		return def
	}
	return g.unknown
}
