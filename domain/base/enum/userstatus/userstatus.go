package userstatus

import (
	"omono/internal/types"
)

// enum for users
const (
	Active    types.Enum = "active"
	Inactive  types.Enum = "inactive"
	Terminate types.Enum = "terminate"
)

// List is used for validation
var List = []types.Enum{
	Active,
	Inactive,
	Terminate,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
