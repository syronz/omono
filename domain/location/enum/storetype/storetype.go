package storetype

import (
	"omono/internal/types"
)

const (
	ShowRoom       types.Enum = "show-room"
	Office         types.Enum = "office"
	HeadQuarter    types.Enum = "head-quarter"
	Branch         types.Enum = "branch"
	Warehosue      types.Enum = "warehouse"
	Representative types.Enum = "representative"
)

var List = []types.Enum{
	ShowRoom,
	Office,
	HeadQuarter,
	Branch,
	Warehosue,
	Representative,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
