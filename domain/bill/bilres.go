package bill

import "omono/internal/types"

// list of resources for bill domain
const (
	Domain string = "bill"

	InvoiceWrite types.Resource = "invoice:write"
	InvoiceRead  types.Resource = "invoice:read"
	InvoiceExcel types.Resource = "invoice:excel"
)
