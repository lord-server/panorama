package game

import (
	"encoding/json"
	"image"
	"os"

	"github.com/weqqr/panorama/mesh"
)

type gameDescriptor struct {
	Aliases map[string]string         `json:"aliases"`
	Nodes   map[string]NodeDescriptor `json:"nodes"`
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

func ResolveNode(descriptor NodeDescriptor, mediaCache *MediaCache) Node {
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

	var descriptor gameDescriptor
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
		node := ResolveNode(gameNode, mediaCache)

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
