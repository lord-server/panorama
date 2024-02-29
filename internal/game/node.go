package game

import (
	"encoding/json"
	"fmt"
)

type DrawType int

const (
	DrawTypeNormal DrawType = iota
	DrawTypeAirlike
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
	DrawTypeNodeBox
	DrawTypeMesh
	DrawTypePlantlikeRooted
)

var DrawTypeNames = map[string]DrawType{
	"node":                      DrawTypeNormal,
	"normal":                    DrawTypeNormal,
	"airlike":                   DrawTypeAirlike,
	"liquid":                    DrawTypeLiquid,
	"flowingliquid":             DrawTypeFlowingLiquid,
	"glasslike":                 DrawTypeGlasslike,
	"glasslike_framed":          DrawTypeGlasslikeFramed,
	"glasslike_framed_optional": DrawTypeGlasslikeFramed,
	"allfaces":                  DrawTypeAllFaces,
	"allfaces_optional":         DrawTypeAllFaces,
	"torchlike":                 DrawTypeTorchlike,
	"signlike":                  DrawTypeSignlike,
	//	"plantlike":                 DrawTypePlantlike,
	"plantlike":        DrawTypeAllFaces,
	"firelike":         DrawTypeFirelike,
	"fencelike":        DrawTypeFencelike,
	"raillike":         DrawTypeRaillike,
	"nodebox":          DrawTypeNodeBox,
	"mesh":             DrawTypeMesh,
	"plantlike_rooted": DrawTypePlantlikeRooted,
}

func (t DrawType) IsLiquid() bool {
	return t == DrawTypeLiquid || t == DrawTypeFlowingLiquid
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

type ParamType int

const (
	ParamTypeLight = iota
	ParamTypeNone
)

var ParamTypeNames = map[string]ParamType{
	"none":  ParamTypeNone,
	"light": ParamTypeLight,
}

func (t *ParamType) UnmarshalJSON(data []byte) error {
	name := "light"
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}

	if paramtype, ok := ParamTypeNames[name]; ok {
		*t = paramtype
	} else {
		return fmt.Errorf("invalid paramtype: `%s`", name)
	}

	return nil
}

type ParamType2 int

const (
	ParamType2FlowingLiquid = iota
	ParamType2WallMounted
	ParamType2FaceDir
	ParamType2Leveled
	ParamType2DegRotate
	ParamType2MeshOptions
	ParamType2Color
	ParamType2ColorFaceDir
	ParamType2ColorWallMounted
	ParamType2GlassLikeLiquidLevel
	ParamType2ColorDegRotate
	ParamType2None
	ParamType2Waving
)

var ParamType2Names = map[string]ParamType2{
	"flowingliquid":        ParamType2FlowingLiquid,
	"wallmounted":          ParamType2WallMounted,
	"facedir":              ParamType2FaceDir,
	"leveled":              ParamType2Leveled,
	"degrotate":            ParamType2DegRotate,
	"meshoptions":          ParamType2MeshOptions,
	"color":                ParamType2Color,
	"colorfacedir":         ParamType2ColorFaceDir,
	"colorwallmounted":     ParamType2ColorWallMounted,
	"glasslikeliquidlevel": ParamType2GlassLikeLiquidLevel,
	"colordegrotate":       ParamType2ColorDegRotate,
	"waving":               ParamType2Waving,
	"none":                 ParamType2None,
}

func (t *ParamType2) UnmarshalJSON(data []byte) error {
	name := "none"
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}

	if paramtype2, ok := ParamType2Names[name]; ok {
		*t = paramtype2
	} else {
		return fmt.Errorf("invalid paramtype2: `%s`", name)
	}

	return nil
}

type NodeBox struct {
	Type  string
	Fixed [][]float64
}

func (n *NodeBox) UnmarshalJSON(data []byte) error {
	type nodeBox struct {
		Type  string        `json:"type"`
		Fixed []interface{} `json:"fixed"`
	}
	inner := &nodeBox{}
	if err := json.Unmarshal(data, inner); err != nil {
		return err
	}

	n.Type = inner.Type
	n.Fixed = make([][]float64, 0)
	if inner.Type != "fixed" {
		return nil
	}

	if len(inner.Fixed) == 0 {
		return nil
	}

	if _, ok := inner.Fixed[0].(float64); ok {
		box := make([]float64, 0)
		for i := 0; i < 6; i++ {
			v, _ := inner.Fixed[i].(float64)
			box = append(box, v)
		}
		n.Fixed = append(n.Fixed, box)
	}

	if _, ok := inner.Fixed[0].([]interface{}); ok {
		for _, boxInterface := range inner.Fixed {
			boxFloat64 := boxInterface.([]interface{})
			box := make([]float64, 0)
			for i := 0; i < 6; i++ {
				v, _ := boxFloat64[i].(float64)
				box = append(box, v)
			}
			n.Fixed = append(n.Fixed, box)
		}
	}

	return nil
}

type NodeDescriptor struct {
	DrawType   DrawType   `json:"drawtype"`
	ParamType  ParamType  `json:"paramtype"`
	ParamType2 ParamType2 `json:"paramtype2"`
	Tiles      []string   `json:"tiles"`
	NodeBox    *NodeBox   `json:"node_box"`
	Mesh       *string    `json:"mesh"`
}

func (n *NodeDescriptor) UnmarshalJSON(data []byte) error {
	type nodeDescriptor NodeDescriptor
	inner := &nodeDescriptor{
		DrawType:   DrawTypeNormal,
		Tiles:      []string{},
		ParamType:  ParamTypeLight,
		ParamType2: ParamType2None,
	}

	if err := json.Unmarshal(data, inner); err != nil {
		return err
	}

	*n = NodeDescriptor(*inner)

	return nil
}
