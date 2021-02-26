package invoicestatus

import (
	"omono/internal/types"
)

const (
	New     types.Enum = "new"
	Pending types.Enum = "pending"
	Done    types.Enum = "done"
	Cancel  types.Enum = "cancel"
	Lock    types.Enum = "lock"
)

var List = []types.Enum{
	New,
	Pending,
	Done,
	Cancel,
	Lock,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
