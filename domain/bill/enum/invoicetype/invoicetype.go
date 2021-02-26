package invoicetype

import (
	"omono/internal/types"
)

const (
	Sale     types.Enum = "sale"
	Purchase types.Enum = "purchase"
	Transfer types.Enum = "transfer"
)

var List = []types.Enum{
	Sale,
	Purchase,
	Transfer,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
