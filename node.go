package main

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

type NodeDef struct {
	drawtype DrawType
	mesh     *Mesh
}
