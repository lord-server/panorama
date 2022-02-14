package game

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"

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
	Tiles    []*image.NRGBA
	Mesh     *mesh.Mesh
}

type Game struct {
	Aliases map[string]string
	Nodes   map[string]Node
	unknown Node
}

func ResolveNode(descriptor NodeDescriptor, mediaCache *MediaCache) Node {
	tiles := make([]*image.NRGBA, len(descriptor.Tiles))
	for i, tileName := range descriptor.Tiles {
		// FIXME: resolve modifiers
		baseImageName := strings.Split(tileName, "^")[0]

		tiles[i] = mediaCache.Image(baseImageName)
	}

	var nodeMesh *mesh.Mesh

	switch descriptor.DrawType {
	case DrawTypeNormal, DrawTypeAllFaces, DrawTypeLiquid, DrawTypeFlowingLiquid, DrawTypeGlasslike, DrawTypeGlasslikeFramed:
		nodeMesh = mesh.Cube()
	case DrawTypeMesh:
		if descriptor.Mesh != nil {
			nodeMesh = mediaCache.Mesh(*descriptor.Mesh)
		}
	}

	return Node{
		DrawType: descriptor.DrawType,
		Tiles:    tiles,
		Mesh:     nodeMesh,
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
		node := ResolveNode(gameNode, mediaCache)

		nodes[name] = node
	}

	return Game{
		Aliases: descriptor.Aliases,
		Nodes:   nodes,
		unknown: Node{
			DrawType: DrawTypeNormal,
			Tiles:    []*image.NRGBA{mediaCache.dummyImage},
			Mesh:     nil,
		},
	}, nil
}

func (g *Game) Node(node string) Node {
	if def, ok := g.Nodes[node]; ok {
		return def
	}
	return g.unknown
}
