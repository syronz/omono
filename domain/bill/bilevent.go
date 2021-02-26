package bill

import "omono/internal/types"

// types for bill domain
const (
	CreateInvoice types.Event = "invoice-create"
	UpdateInvoice types.Event = "invoice-update"
	DeleteInvoice types.Event = "invoice-delete"
	ListInvoice   types.Event = "invoice-list"
	ViewInvoice   types.Event = "invoice-view"
	ExcelInvoice  types.Event = "invoice-excel"
)
