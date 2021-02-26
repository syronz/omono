package pricemode

import (
	"omono/internal/types"
)

const (
	Whole       types.Enum = "whole"
	VIP         types.Enum = "vip"
	Distributor types.Enum = "distributor"
	Export      types.Enum = "export"
	Retail      types.Enum = "retail"
	Old         types.Enum = "old"
)

var List = []types.Enum{
	Whole,
	VIP,
	Distributor,
	Export,
	Retail,
	Old,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
