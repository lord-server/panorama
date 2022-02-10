package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

type DrawType int

const (
	DrawTypeNormal DrawType = iota
	DrawTypeAirLlke
	DrawTypeLiquid
	DrawTypeFlowingLiquid
	DrawTypeGlasslike
	DrawTypeGlasslikeFramed
	DrawTypeAllFaces
	DrawTypeTorchlike
	DrawTypeSignlike
	DrawTypePlantlike
	DrawTypeFirelike
	DrawTypeFencelike
	DrawTypeRaillike
	DrawTypeNodebox
	DrawTypeMesh
	DrawTypePlantlikeRooted
)

var DrawTypeNames = map[string]DrawType{
	"normal":                    DrawTypeNormal,
	"airlike":                   DrawTypeAirLlke,
	"liquid":                    DrawTypeLiquid,
	"flowingliquid":             DrawTypeFlowingLiquid,
	"glasslike":                 DrawTypeGlasslike,
	"glasslike_framed":          DrawTypeGlasslikeFramed,
	"glasslike_framed_optional": DrawTypeGlasslikeFramed,
	"allfaces":                  DrawTypeAllFaces,
	"allfaces_optional":         DrawTypeAllFaces,
	"torchlike":                 DrawTypeTorchlike,
	"signlike":                  DrawTypeSignlike,
	"plantlike":                 DrawTypePlantlike,
	"firelike":                  DrawTypeFirelike,
	"fencelike":                 DrawTypeFencelike,
	"raillike":                  DrawTypeRaillike,
	"nodebox":                   DrawTypeNodebox,
	"mesh":                      DrawTypeMesh,
	"plantlike_rooted":          DrawTypePlantlikeRooted,
}

func (t *DrawType) UnmarshalJSON(data []byte) error {
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}

	if drawtype, ok := DrawTypeNames[name]; ok {
		*t = drawtype
	} else {
		return fmt.Errorf("invalid drawtype: `%s`", name)
	}

	return nil
}

type GameNodeDef struct {
	DrawType DrawType `json:"drawtype"`
	Tiles    []string `json:"tiles"`
	Mesh     *string  `json:"mesh"`
}

func (n *GameNodeDef) UnmarshalJSON(data []byte) error {
	type DefaultNodeDef GameNodeDef
	inner := &DefaultNodeDef{
		DrawType: DrawTypeNormal,
		Tiles:    []string{},
	}

	if err := json.Unmarshal(data, inner); err != nil {
		return err
	}

	*n = GameNodeDef(*inner)
	return nil
}

type GameDefinitions struct {
	Aliases map[string]string      `json:"aliases"`
	Nodes   map[string]GameNodeDef `json:"nodes"`
}

func toNRGBA(img image.Image) *image.NRGBA {
	dst := image.NewNRGBA(img.Bounds())
	draw.Draw(dst, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return dst
}

func readPNG(path string) (*image.NRGBA, error) {
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

type NodeDef struct {
	DrawType DrawType
	Tiles    []*image.NRGBA
	Mesh     *Mesh
}

func LoadNodeDef(nd GameNodeDef, mediaCache *MediaCache) NodeDef {
	return NodeDef{
		DrawType: nd.DrawType,
	}
}

type Game struct {
	aliases map[string]string
	nodes   map[string]NodeDef
	unknown NodeDef
}

func ResolveNodeDef(gameNode GameNodeDef, mediaCache *MediaCache) NodeDef {
	tiles := make([]*image.NRGBA, len(gameNode.Tiles))
	for i, tileName := range gameNode.Tiles {
		// FIXME: resolve modifiers
		baseImageName := strings.Split(tileName, "^")[0]

		tiles[i] = mediaCache.Image(baseImageName)
	}

	var mesh *Mesh
	if gameNode.DrawType == DrawTypeNormal {
		mesh = Cube()
	}

	if gameNode.Mesh != nil {
		mesh = mediaCache.meshes[*gameNode.Mesh]
	}

	return NodeDef{
		DrawType: gameNode.DrawType,
		Tiles:    tiles,
		Mesh:     mesh,
	}
}

func LoadGame(desc string, path string) (Game, error) {
	descJSON, err := os.ReadFile(desc)
	if err != nil {
		return Game{}, err
	}

	var defs GameDefinitions
	err = json.Unmarshal(descJSON, &defs)
	if err != nil {
		return Game{}, err
	}

	mediaCache := NewMediaCache()

	err = mediaCache.fetchMedia(path)
	if err != nil {
		return Game{}, err
	}

	nodes := make(map[string]NodeDef)
	for name, gameNode := range defs.Nodes {
		node := ResolveNodeDef(gameNode, mediaCache)

		nodes[name] = node
	}

	return Game{
		aliases: defs.Aliases,
		nodes:   nodes,
		unknown: NodeDef{
			DrawType: DrawTypeNormal,
			Tiles:    []*image.NRGBA{mediaCache.dummyImage},
			Mesh:     nil,
		},
	}, nil
}

func (g *Game) NodeDef(node string) NodeDef {
	if def, ok := g.nodes[node]; ok {
		return def
	}
	return g.unknown
}
