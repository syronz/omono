package notification

import "omono/internal/types"

// list of resources for notification domain
const (
	Domain string = "notification"

	MessageWrite types.Resource = "message:write"
	MessageRead  types.Resource = "message:read"
	MessageExcel types.Resource = "message:excel"
)
