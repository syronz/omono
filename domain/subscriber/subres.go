package subscribe

import "omono/internal/types"

// list of resources for subscribe domain
const (
	Domain string = "subscribe"

	AccountRead  types.Resource = "account:read"
	AccountWrite types.Resource = "account:write"
	AccountExcel types.Resource = "account:excel"

	PhoneRead  types.Resource = "phone:read"
	PhoneWrite types.Resource = "phone:write"
	PhoneExcel types.Resource = "phone:excel"
)
