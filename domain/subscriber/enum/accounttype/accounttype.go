package accounttype

import (
	"omono/internal/types"
)

//the type of accounts
const (
	Regular    types.Enum = "regular"
	Business   types.Enum = "business"
	Enterprise types.Enum = "enterprise"
	VIP        types.Enum = "vip"
	Employee   types.Enum = "employee"
	Government types.Enum = "government"
)

// List used for validation
var List = []types.Enum{
	Regular,
	Business,
	Enterprise,
	VIP,
	Employee,
	Government,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
