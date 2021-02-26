package transactionstatus

import "omono/internal/types"

// Transactions type
const (
	Approved   types.Enum = "approved"
	Unapproved types.Enum = "unapproved"
)

//List of transaction status
var List = []types.Enum{
	Approved,
	Unapproved,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
