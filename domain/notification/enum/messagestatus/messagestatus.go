package messagestatus

import "omono/internal/types"

// Transactions type
const (
	New  types.Enum = "new"
	Seen types.Enum = "seen"
	Sent types.Enum = "sent"
)

//List of message status
var List = []types.Enum{
	New,
	Seen,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}
