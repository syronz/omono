package storestatus

import (
	"omono/internal/types"
)

const (
	Active   types.Enum = "active"
	Inactive types.Enum = "inactive"
)

var List = []types.Enum{
	Active,
	Inactive,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
