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
	"plantlike":                 DrawTypePlantlike,
	"firelike":                  DrawTypeFirelike,
	"fencelike":                 DrawTypeFencelike,
	"raillike":                  DrawTypeRaillike,
	"nodebox":                   DrawTypeNodeBox,
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
	var name string
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

type NodeDescriptor struct {
	DrawType  DrawType  `json:"drawtype"`
	ParamType ParamType `json:"paramtype"`
	Tiles     []string  `json:"tiles"`
	Mesh      *string   `json:"mesh"`
}

func (n *NodeDescriptor) UnmarshalJSON(data []byte) error {
	type nodeDescriptor NodeDescriptor
	inner := &nodeDescriptor{
		DrawType: DrawTypeNormal,
		Tiles:    []string{},
	}

	if err := json.Unmarshal(data, inner); err != nil {
		return err
	}

	*n = NodeDescriptor(*inner)

	return nil
}
