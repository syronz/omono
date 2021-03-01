package segment

import "omono/internal/types"

// list of resources for subscribe domain
const (
	Domain string = "subscribe"

	CompanyRead  types.Resource = "company:read"
	CompanyWrite types.Resource = "company:write"
	CompanyExcel types.Resource = "company:excel"
)
