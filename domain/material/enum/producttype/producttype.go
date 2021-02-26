package producttype

import (
	"omono/internal/types"
)

const (
	Serial   types.Enum = "serial"
	Stocking types.Enum = "stocking"
	Service  types.Enum = "service"
)

var List = []types.Enum{
	Serial,
	Stocking,
	Service,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
